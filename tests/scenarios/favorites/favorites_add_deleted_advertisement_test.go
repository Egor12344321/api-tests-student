package favorites

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/favorites"
)

func (s *TestSuite) TestAddFavoriteDeletedAdvertisement() {
	const expectedStatusCode = http.StatusNotFound

	var advertisementID string

	if !s.Run("Владелец создаёт и удаляет объявление", func() {
		advertisement := createAdvertisementForFavorites(s.T(), s.ownerLoginData.Token)
		advertisementID = advertisement.ID
		// фикстура удалит объявление через t.Cleanup в конце этого шага.
	}) {
		return
	}

	s.Run("Покупатель пытается добавить удалённое объявление в избранное", func() {
		t := s.T()
		t.Logf("Добавляем удалённое объявление %s, ожидаем HTTP %d", advertisementID, expectedStatusCode)

		errorResponse := favorites.AddFavoriteExpectError(t, s.buyerLoginData.Token, advertisementID, expectedStatusCode)

		assertions.CheckErrorResponse(t, errorResponse, "not_found", "Advertisement not found")
	})
}
