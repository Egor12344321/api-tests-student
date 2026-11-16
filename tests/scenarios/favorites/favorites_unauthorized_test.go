package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	favModels "api-tests-template/internal/managers/favorites/models"

	base "api-tests-template/tests"
)


// параметризированный тест для проверки запросов без авторизации
func TestFavoritesUnauthorized(t *testing.T) {
	const (
		advertisementID = "UnrealEgorsIdAdv"
		incorrectToken  = "incorrect_token"
		errorCode       = "unauthorized"
	)

	base.SetupSuite()
	base.Precondition("Проверяем избранное без действующего токена")

	endpoints := []struct {
		name     string
		argument string
		request  func(t *testing.T, token, argument string, expectedStatusCode int) favModels.ErrorResponse
	}{
		{"GET /favorites", "", favorites.GetFavoritesExpectError},
		{"POST /favorites", advertisementID, favorites.AddFavoriteExpectError},
		{"DELETE /favorites/{ad_id}", advertisementID, favorites.DeleteFavoriteExpectError},
	}

	cases := []struct {
		name    string
		token   string
		message string
	}{
		{"Без токена", "", "Authorization header required"},
		{"С неверным токеном", incorrectToken, "Invalid or expired token"},
	}

	for _, endpoint := range endpoints {
		for _, tc := range cases {
			t.Run(endpoint.name+" — "+tc.name, func(t *testing.T) {
				t.Log("Отправляем запрос и ожидаем HTTP 401")
				errorResponse := endpoint.request(t, tc.token, endpoint.argument, http.StatusUnauthorized)

				require.Equal(t, errorCode, errorResponse.Error, "Код ошибки не совпадает")
				require.Equal(t, tc.message, errorResponse.Message, "Сообщение об ошибке не совпадает")
			})
		}
	}
}