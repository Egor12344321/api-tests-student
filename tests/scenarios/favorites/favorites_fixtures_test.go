package favorites

import (
	"testing"

	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/utils"
)

// createAdvertisementForFavorites создаёт объявление и регистрирует его удаление после теста.
func createAdvertisementForFavorites(t *testing.T, ownerToken string) advModels.Advertisement {
	t.Helper()

	title := "Избранное " + utils.RandomString(12)
	fields := map[string]string{
		"title":       title,
		"description": "Объявление для проверки избранного",
		"price":       "1000",
		"quantity":    "3",
	}

	created := advertisements.CreateAdvertisement(t, ownerToken, fields, []string{photoPath})
	t.Cleanup(func() {
		t.Logf("Удаляем тестовое объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, ownerToken, created.ID)
	})

	return advModels.Advertisement{
		ID:    created.ID,
		Title: title,
	}
}
