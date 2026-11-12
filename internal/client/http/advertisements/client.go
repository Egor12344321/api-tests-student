package advertisements

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// запрос на получение объявления по названию
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).
		Status(http.StatusOK).
		End().Response
}

// запрос на создание объявления
func HttpPostAdvertisement(t *testing.T, token string, fields map[string]string, photoPaths []string) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Post(path.AdvertisementPath)
	response := request.Expect(t)

	for name, value := range fields {
		request.MultipartFormData(name, value)
	}
	request.MultipartFile("photos", photoPaths...)

	return response.End().Response
}