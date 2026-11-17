package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
	favModels "api-tests-template/internal/managers/favorites/models"

	base "api-tests-template/tests"
)

// TestFavoritesUnauthorized проверяет доступ к избранному без действующего токена.
func TestFavoritesUnauthorized(t *testing.T) {
	const (
		advertisementID    = "39960e58-db41-41ba-959e-01073a913692"
		incorrectToken     = "incorrect_token"
		errorCode          = "unauthorized"
		expectedStatusCode = http.StatusUnauthorized
	)

	base.SetupSuite()
	t.Logf("Для всех запросов ожидаем HTTP %d", expectedStatusCode)

	endpoints := []struct {
		name     string
		argument string
		request  func(t *testing.T, token, argument string, expectedStatusCode int) favModels.ErrorResponse
	}{
		{"Запрашиваем список избранного", "", favorites.GetFavoritesExpectError},
		{"Добавляем объявление в избранное", advertisementID, favorites.AddFavoriteExpectError},
		{"Удаляем объявление из избранного", advertisementID, favorites.DeleteFavoriteExpectError},
	}

	cases := []struct {
		name    string
		token   string
		message string
	}{
		{"без токена", "", "Authorization header required"},
		{"с неверным токеном", incorrectToken, "Invalid or expired token"},
	}

	for _, endpoint := range endpoints {
		for _, tc := range cases {
			t.Run(endpoint.name+" "+tc.name, func(t *testing.T) {
				errorResponse := endpoint.request(t, tc.token, endpoint.argument, expectedStatusCode)

				assertions.CheckErrorResponse(t, errorResponse, errorCode, tc.message)
			})
		}
	}
}
