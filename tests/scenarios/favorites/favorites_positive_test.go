package favorites

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const photoPath = "../advertisementCreate/fixtures/photoForCreateAdv.png"

type TestSuite struct {
	suite.Suite
	ownerLoginData models.LoginOkResponse
	buyerLoginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	t := s.T()
	buyerLogin := os.Getenv("TEST_BUYER_LOGIN")
	buyerPassword := os.Getenv("TEST_BUYER_PASSWORD")
	require.NotEmpty(t, buyerLogin, "задайте TEST_BUYER_LOGIN в .env.override")
	require.NotEmpty(t, buyerPassword, "задайте TEST_BUYER_PASSWORD в .env.override")

	base.Precondition("Авторизуем владельца объявления и покупателя")
	s.ownerLoginData = auth.Login(t, os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
	s.buyerLoginData = auth.Login(t, buyerLogin, buyerPassword)
	require.NotEqual(t, s.ownerLoginData.User.Id, s.buyerLoginData.User.Id, "Владелец и покупатель должны быть разными пользователями")
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestFavoritesPositive() {
	t := s.T()
	ownerToken := s.ownerLoginData.Token
	buyerToken := s.buyerLoginData.Token
	title := "Избранное " + utils.RandomString(12)
	fields := map[string]string{
		"title":       title,
		"description": "Объявление для проверки избранного",
		"price":       "1000",
		"quantity":    "3",
	}

	t.Log("Владелец создаёт объявление")
	created := advertisements.CreateAdvertisement(t, ownerToken, fields, []string{photoPath})
	advertisementDeleted := false
	t.Cleanup(func() {
		if !advertisementDeleted {
			t.Logf("Удаляем тестовое объявление %s", created.ID)
			advertisements.DeleteAdvertisement(t, ownerToken, created.ID)
		}
	})

	if !s.Run("Добавляем объявление в избранное и находим его в списке", func() {
		t := s.T()
		favorites.AddFavorite(t, buyerToken, created.ID)
		require.True(t, advertisementInFavorites(t, buyerToken, title, created.ID), "Объявление %s не появилось в избранном", created.ID)
	}) {
		return
	}

	if !s.Run("Удаляем объявление из избранного и проверяем список", func() {
		t := s.T()
		favorites.DeleteFavorite(t, buyerToken, created.ID)
		require.False(t, advertisementInFavorites(t, buyerToken, title, created.ID), "Объявление %s осталось в списке после удаления из избранного", created.ID)
	}) {
		return
	}

	s.Run("Удаляем само объявление и проверяем список избранного", func() {
		t := s.T()
		favorites.AddFavorite(t, buyerToken, created.ID)
		require.True(t, advertisementInFavorites(t, buyerToken, title, created.ID), "Объявление %s не появилось после повторного добавления", created.ID)

		advertisements.DeleteAdvertisement(t, ownerToken, created.ID)
		advertisementDeleted = true

		require.False(t, advertisementInFavorites(t, buyerToken, title, created.ID), "Объявление %s осталось в избранном после удаления самого объявления", created.ID)
	})
}

// advertisementInFavorites проверяет наличие объявления по ID в списке избранного.
func advertisementInFavorites(t *testing.T, token, title, advertisementID string) bool {
	t.Helper()

	favoriteListResponse := favorites.GetFavorites(t, token, title)
	for _, item := range favoriteListResponse.Items {
		if item.Advertisement.ID == advertisementID {
			return true
		}
	}

	return false
}
