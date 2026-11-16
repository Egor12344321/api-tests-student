package favorites

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/favorites/models"
)

// AddFavorite добавляет объявление в избранное и проверяет статус 201.
func AddFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	request := models.AddFavoriteRequest{AdvertisementID: advertisementID}
	raw := favorites.HttpPostFavorite(t, token, request)
	readFavoriteResponse(t, raw, http.StatusCreated)
}

// GetFavorites возвращает список избранного с поиском по названию.
func GetFavorites(t *testing.T, token, textForSearch string) models.FavoriteListResponse {
	t.Helper()

	raw := favorites.HttpGetFavorites(t, token, textForSearch)
	body := readFavoriteResponse(t, raw, http.StatusOK)

	favoriteListResponse := models.FavoriteListResponse{}
	require.NoError(t, json.Unmarshal(body, &favoriteListResponse), "Не удалось разобрать список избранного")

	return favoriteListResponse
}

// DeleteFavorite удаляет объявление из избранного и проверяет статус 204.
func DeleteFavorite(t *testing.T, token, advertisementID string) {
	t.Helper()

	raw := favorites.HttpDeleteFavorite(t, token, advertisementID)
	readFavoriteResponse(t, raw, http.StatusNoContent)
}

// readFavoriteResponse читает и закрывает тело ответа и проверяет HTTP-статус.
func readFavoriteResponse(t *testing.T, raw *http.Response, expectedStatusCode int) []byte {
	t.Helper()
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ избранного")
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "HTTP status code должен быть %d. Тело ответа: %s", expectedStatusCode, body)

	return body
}

// GetFavoritesExpectError проверяет ожидаемый статус при чтении избранного.
func GetFavoritesExpectError(t *testing.T, token, textForSearch string, expectedStatusCode int) models.ErrorResponse {
	t.Helper()

	raw := favorites.HttpGetFavorites(t, token, textForSearch)
	return readFavoriteErrorResponse(t, raw, expectedStatusCode)
}

// AddFavoriteExpectError проверяет ожидаемый статус при добавлении в избранное.
func AddFavoriteExpectError(t *testing.T, token, advertisementID string, expectedStatusCode int) models.ErrorResponse {
	t.Helper()

	request := models.AddFavoriteRequest{AdvertisementID: advertisementID}
	raw := favorites.HttpPostFavorite(t, token, request)
	return readFavoriteErrorResponse(t, raw, expectedStatusCode)
}

// DeleteFavoriteExpectError проверяет ожидаемый статус при удалении из избранного.
func DeleteFavoriteExpectError(t *testing.T, token, advertisementID string, expectedStatusCode int) models.ErrorResponse {
	t.Helper()

	raw := favorites.HttpDeleteFavorite(t, token, advertisementID)
	return readFavoriteErrorResponse(t, raw, expectedStatusCode)
}

// readFavoriteErrorResponse проверяет статус и разбирает тело ошибки.
func readFavoriteErrorResponse(t *testing.T, raw *http.Response, expectedStatusCode int) models.ErrorResponse {
	t.Helper()

	body := readFavoriteResponse(t, raw, expectedStatusCode)
	errorResponse := models.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &errorResponse), "Не удалось разобрать ответ с ошибкой избранного")

	return errorResponse
}