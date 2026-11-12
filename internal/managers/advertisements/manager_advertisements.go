package advertisements

import (
	"encoding/json"
	"io"
	"testing"
	"net/http"

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


func CreateAdvertisement(t *testing.T, token string, fields map[string]string, photoPaths []string) models.AdvertisementCreateResponse {
	raw := advertisements.HttpPostAdvertisement(t, token, fields, photoPaths)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ на создание объявления")
	require.Equalf(t, http.StatusCreated, raw.StatusCode, "Тело ответа %s", body)

	var result models.AdvertisementCreateResponse
	require.NoError(t, json.Unmarshal(body, &result), "Не удалось разобрать ответ на создание объявления")
	require.NotEmpty(t, result.ID, "Сервер не вернул ID созданного объявления")

	return result
}


func GetAdvertisementByID(t *testing.T, id string) models.AdvertisementDetails {
	raw := advertisements.HttpGetAdvertisementByID(t, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ для объявления " + id)
	require.Equalf(t, http.StatusOK, raw.StatusCode, "Тело ответа %s", body)

	var result models.AdvertisementDetails
	require.NoError(t, json.Unmarshal(body, &result), "Не удалось разобрать данные объявления")

	return result
}


func DeleteAdvertisement(t *testing.T, token, id string) {
	raw := advertisements.HttpDeleteAdvertisement(t, token, id)
	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	require.NoError(t, err, "Не удалось прочитать ответ на удаление объявления")
	require.Equalf(t, http.StatusNoContent, raw.StatusCode, "Удаление объявления %s. Тело ответа %s", id, body)
}