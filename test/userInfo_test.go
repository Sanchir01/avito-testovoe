package test

import (
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
)

func TestGetUserInfoE2E(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	e := httpexpect.Default(t, u.String())
	testUser := map[string]interface{}{
		"email":    "test@test.ru",
		"password": "test01",
	}
	authResp := e.
		POST("/api/auth").
		WithJSON(testUser).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object()

	token := authResp.Value("token").String().Raw()

	t.Run("testing info user", func(_ *testing.T) {
		e.GET("/api/info").
			WithHeader("Authorization", "Bearer "+token).
			Expect().
			Status(http.StatusOK)
	})

	t.Run("Unauthorized", func(_ *testing.T) {
		e.GET("/api/info").
			Expect().
			Status(http.StatusUnauthorized)
	})
}
