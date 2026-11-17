package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
)

// TestAddFavoriteOwnAdvertisement проверяет запрет добавления своего объявления в избранное.
func (s *TestSuite) TestAddFavoriteOwnAdvertisement() {
	const expectedStatusCode = http.StatusForbidden

	t := s.T()
	ownerToken := s.ownerLoginData.Token

	t.Logf("Создаём своё объявление. При добавлении в избранное ожидаем HTTP %d", expectedStatusCode)
	advertisement := createAdvertisementForFavorites(t, ownerToken)

	if !s.Run("Пытаемся добавить своё объявление в избранное", func() {
		t := s.T()
		errorResponse := favorites.AddFavoriteExpectError(t, ownerToken, advertisement.ID, expectedStatusCode)

		assertions.CheckErrorResponse(t, errorResponse, "forbidden", "Cannot add own advertisement to favorites")
	}) {
		return
	}

	s.Run("Проверяем, что своё объявление не добавилось в избранное", func() {
		t := s.T()
		require.False(t, advertisementInFavorites(t, ownerToken, advertisement.Title, advertisement.ID), "Собственное объявление %s появилось в избранном", advertisement.ID)
	})
}
