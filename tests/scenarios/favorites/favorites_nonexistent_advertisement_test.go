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
		advertisementID    = "39960e58-db41-41ba-959e-01073a913692"
		errorCode          = "not_found"
		expectedStatusCode = http.StatusNotFound
	)

	s.T().Logf(
		"Проверяем попытки добавить несуществующее объявление в избранное и удалить его из избранного. ID %s, ожидаем HTTP %d",
		advertisementID, expectedStatusCode,
	)

	cases := []struct {
		name    string
		request func(t *testing.T, token, advertisementID string, expectedStatusCode int) favModels.ErrorResponse
		message string
	}{
		{"Пытаемся добавить несуществующее объявление в избранное", favorites.AddFavoriteExpectError, "Advertisement not found"},
		{"Пытаемся удалить несуществующее объявление из избранного", favorites.DeleteFavoriteExpectError, "Favorite entry not found"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			errorResponse := tc.request(t, s.buyerLoginData.Token, advertisementID, expectedStatusCode)

			require.Equal(t, errorCode, errorResponse.Error, "Некорректное значение error параметра")
			require.Equal(t, tc.message, errorResponse.Message, "Некорректное значение message параметра")
		})
	}
}