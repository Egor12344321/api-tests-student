package advertisements

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
)

func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	raw := advertisements.HttpGetAdvertisements(t, token, text)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err)

	advertisementsResponse := models.AdvertisementsList{}
	require.NoError(t, json.Unmarshal(body, &advertisementsResponse))

	return advertisementsResponse
}

// CreateAdvertisement создаёт объявление и возвращает его ID и фотографии.
func CreateAdvertisement(t *testing.T, token string, fields map[string]string, photoPaths []string) models.AdvertisementCreateResponse {
	body := postAdvertisement(t, token, fields, photoPaths, http.StatusCreated)

	advertisementCreateResponse := models.AdvertisementCreateResponse{}
	require.NoError(t, json.Unmarshal(body, &advertisementCreateResponse), "Не удалось разобрать ответ на создание объявления")
	require.NotEmpty(t, advertisementCreateResponse.ID, "Сервер не вернул ID созданного объявления")

	return advertisementCreateResponse
}

// CreateAdvertisementExpectError проверяет ожидаемый статус и возвращает ошибку API.
func CreateAdvertisementExpectError(t *testing.T, token string, fields map[string]string, photoPaths []string, expectedStatusCode int) models.ErrorResponse {
	body := postAdvertisement(t, token, fields, photoPaths, expectedStatusCode)

	errorResponse := models.ErrorResponse{}
	require.NoError(t, json.Unmarshal(body, &errorResponse), "Не удалось разобрать ответ с ошибкой")

	return errorResponse
}

// postAdvertisement отправляет запрос на создание, читает тело и проверяет статус (вынес, чтобы переиспользовать код).
func postAdvertisement(t *testing.T, token string, fields map[string]string, photoPaths []string, expectedStatusCode int) []byte {
	t.Helper()

	raw := advertisements.HttpPostAdvertisement(t, token, fields, photoPaths)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ на создание объявления")
	require.Equalf(t, expectedStatusCode, raw.StatusCode, "Тело ответа: %s", body)

	return body
}

// GetAdvertisementByID возвращает данные объявления по его ID.
func GetAdvertisementByID(t *testing.T, id string) models.AdvertisementDetails {
	raw := advertisements.HttpGetAdvertisementByID(t, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ для объявления "+id)
	require.Equalf(t, http.StatusOK, raw.StatusCode, "Тело ответа %s", body)

	var advertisementDetailsResponse models.AdvertisementDetails
	require.NoError(t, json.Unmarshal(body, &advertisementDetailsResponse), "Не удалось разобрать данные объявления")

	return advertisementDetailsResponse
}

// DeleteAdvertisement удаляет объявление по его ID и проверяет статус ответа.
func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ на удаление объявления")
	require.Equalf(t, http.StatusNoContent, raw.StatusCode, "Удаление объявления %s. Тело ответа %s", id, body)
}

// GetAdvertisementPhotos возвращает список фотографий объявления.
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.PhotoResponse {
	raw := advertisements.HttpGetAdvertisementPhotos(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ с фотографиями")
	require.Equalf(t, http.StatusOK, raw.StatusCode, "Тело ответа %s", body)

	var photoResponse []models.PhotoResponse
	require.NoError(t, json.Unmarshal(body, &photoResponse), "Не удалось разобрать список фотографий")

	return photoResponse
}
