package favorites

import (
	"fmt"
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// HttpPostFavorite добавляет объявление в избранное пользователя.
func HttpPostFavorite(t *testing.T, token string, request models.AddFavoriteRequest) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpGetFavorites получает избранные объявления с поиском по названию.
func HttpGetFavorites(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Query("search", textForSearch).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного по ID объявления.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(fmt.Sprintf(path.FavoriteAdvertisementPath, advertisementID)).
		Expect(t).
		End().Response
}
