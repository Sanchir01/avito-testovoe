package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
	"github.com/Sanchir01/avito-testovoe/internal/feature/user/mocks"
	"github.com/Sanchir01/avito-testovoe/pkg/lib/api"
	sl "github.com/Sanchir01/avito-testovoe/pkg/lib/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	cases := []struct {
		email     string
		password  string
		respError string
		mockToken string
		mockError error
	}{
		{
			email:     "test@test.ru",
			password:  "test01",
			mockToken: "token1",
		},
		{
			email:     "test1@test.ru",
			password:  "test02",
			mockToken: "token2",
		},
		{
			email:     "test2@test.ru",
			password:  "test03",
			mockToken: "token3",
		},
		{
			email:     "test@test.ru",
			password:  "test02",
			mockToken: "",
			respError: "Введен неправильный пароль",
			mockError: api.ErrWrongPasswordError,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.email, func(t *testing.T) {
			t.Parallel()
			userservice := mocks.NewHandlerUser(t)

			userservice.On("Auth", mock.Anything, tc.email, tc.password).Return(tc.mockToken, tc.mockError).Once()
			handler := user.NewHandler(userservice, sl.NewDiscardLogger())
			input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tc.email, tc.password)

			req, err := http.NewRequest(http.MethodPost, "/api/auth", bytes.NewReader([]byte(input)))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()

			handler.AuthHandler(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
			body := rr.Body.String()

			var resp user.AuthResponse

			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}
}
