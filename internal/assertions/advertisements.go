package assertions

import (
	"testing"

	"api-tests-template/internal/managers/advertisements/models"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()

	for _, item := range items {
		if item.Title == title {
			return
		}
	}

	t.Fatalf("Нет объявления с названием %q", title)
}
