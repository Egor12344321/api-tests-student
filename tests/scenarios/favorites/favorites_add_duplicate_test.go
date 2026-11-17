package favorites

import (
	"net/http"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
)

// TestAddFavoriteDuplicate проверяет запрет повторного добавления объявления в избранное.
func (s *TestSuite) TestAddFavoriteDuplicate() {
	const expectedStatusCode = http.StatusConflict

	t := s.T()
	buyerToken := s.buyerLoginData.Token

	t.Log("Создаём чужое объявление и добавляем его в избранное покупателя")
	advertisement := createAdvertisementForFavorites(t, s.ownerLoginData.Token)
	favorites.AddFavorite(t, buyerToken, advertisement.ID)

	t.Logf("При повторном добавлении ожидаем HTTP %d", expectedStatusCode)
	if !s.Run("Пытаемся повторно добавить объявление в избранное", func() {
		t := s.T()
		errorResponse := favorites.AddFavoriteExpectError(t, buyerToken, advertisement.ID, expectedStatusCode)

		assertions.CheckErrorResponse(t, errorResponse, "already_exists", "Advertisement is already in favorites")
	}) {
		return
	}

	s.Run("Проверяем, что объявление осталось в избранном в одном экземпляре", func() {
		t := s.T()
		favoriteListResponse := favorites.GetFavorites(t, buyerToken, advertisement.Title)

		require.Len(t, favoriteListResponse.Items, 1, "В списке должна быть одна запись")
		require.Equal(t, advertisement.ID, favoriteListResponse.Items[0].Advertisement.ID, "ID объявления не совпадает")
	})
}
