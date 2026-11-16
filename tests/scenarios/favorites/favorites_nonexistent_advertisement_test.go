package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/favorites"
	favModels "api-tests-template/internal/managers/favorites/models"
)

func (s *TestSuite) TestFavoritesNonexistentAdvertisement() {
	const (
		advertisementID = "39960e58-db41-41ba-959e-01073a913692"
		errorCode       = "not_found"
	)

	cases := []struct {
		name    string
		request func(t *testing.T, token, advertisementID string, expectedStatusCode int) favModels.ErrorResponse
		message string
	}{
		{"Добавляем несуществующее объявление в избранное", favorites.AddFavoriteExpectError, "Advertisement not found"},
		{"Удаляем несуществующее объявление из избранного", favorites.DeleteFavoriteExpectError, "Favorite entry not found"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			t.Logf("Отправляем запрос с ID %s и ожидаем HTTP 404", advertisementID)
			errorResponse := tc.request(t, s.buyerLoginData.Token, advertisementID, http.StatusNotFound)

			require.Equal(t, errorCode, errorResponse.Error, "Код ошибки не совпадает")
			require.Equal(t, tc.message, errorResponse.Message, "Сообщение об ошибке не совпадает")
		})
	}
}
