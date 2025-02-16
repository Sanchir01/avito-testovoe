package test

import (
	"github.com/Sanchir01/avito-testovoe/internal/feature/user"
	"github.com/gavv/httpexpect/v2"
	"net/http"
	"net/url"
	"testing"
)

func TestAuthE2E(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   "localhost:8080",
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/api/auth").WithJSON(user.AuthRequest{
		Email:    "test@test.com",
		Password: "test01",
	}).
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().
		ContainsKey("token")
}

func TestBuyProductE2E(t *testing.T) {
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

	obj := e.GET("/api/products").
		Expect().
		Status(http.StatusOK).
		JSON().
		Object().ContainsKey("products")
	products := obj.Value("products").Array()

	firstProductID := products.Value(0).Object().Value("ID").String().Raw()

	t.Run("test buy product", func(_ *testing.T) {
		resp := e.GET("/api/buy/{item}", firstProductID).
			WithHeader("Authorization", "Bearer "+token).
			Expect()
		statusCode := resp.Raw().StatusCode

		if statusCode != http.StatusOK {
			errObj := resp.JSON().Object()
			errMsg := errObj.Value("error").String().Raw()
			t.Errorf("Ошибка при покупке продукта. Статус: %d, сообщение: %s", statusCode, errMsg)
		}
	})

}
