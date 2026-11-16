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
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "Тело ответа: %s", body)

	return body
}
