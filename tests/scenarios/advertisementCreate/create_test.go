package advertisementCreate

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/managers/advertisements"
	advModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	description = "Описание тестового объявления"
	price       = 1000
	quantity    = 3
	photoPath   = "fixtures/photoForCreateAdv.png"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// newAdvertisementFields возвращает корректные поля формы, чтобы использовать в негативных сценариях и они не влияли друг на друга, то есть каждый тест получит полностью заполненный adv и изменит какое-то поле под проверку
func newAdvertisementFields() map[string]string {
	return map[string]string{
		"title":       "Тестовое объявление: " + utils.RandomString(12),
		"description": description,
		"price":       strconv.Itoa(price),
		"quantity":    strconv.Itoa(quantity),
	}
}

// проверка позитивного сценария добавления объявления с проверкой поиска по названию и проверкой доступности фото
func (s *TestSuite) TestCreateAdvertisementPositive() {
	const (
		searchAttempts = 5
		searchInterval = 500 * time.Millisecond
	)

	t := s.T()
	token := s.loginData.Token
	fields := newAdvertisementFields()
	title := fields["title"]

	var created advModels.AdvertisementCreateResponse
	if !s.Run("Создаём объявление с фотографией", func() {
		created = advertisements.CreateAdvertisement(s.T(), token, fields, []string{photoPath})
	}) {
		return
	}

	t.Cleanup(func() {
		t.Logf("Удаляем объявление %s", created.ID)
		advertisements.DeleteAdvertisement(t, token, created.ID)
	})

	s.Run("Проверяем сохранённые поля объявления", func() {
		saved := advertisements.GetAdvertisementByID(s.T(), created.ID)

		require.Equal(s.T(), created.ID, saved.ID, "ID объявления не совпадает")
		require.Equal(s.T(), title, saved.Title, "Название не совпадает")
		require.Equal(s.T(), description, saved.Description, "Описание не совпадает")
		require.Equal(s.T(), float64(price), saved.Price, "Цена не совпадает")
		require.Equal(s.T(), quantity, saved.Quantity, "Количество не совпадает")
	})

	// проверка появления объявления в поиске с пуллингом, чтобы объявление успело проиндексироваться, если используется соответствующая БД
	s.Run("Находим объявление по полному названию", func() {
		for attempt := 1; attempt <= searchAttempts; attempt++ {
			list := advertisements.GetAdvertisementsList(s.T(), token, title)
			for _, item := range list.Items {
				if item.ID == created.ID && item.Title == title {
					return
				}
			}

			if attempt < searchAttempts {
				time.Sleep(searchInterval)
			}
		}

		require.True(s.T(), false, "Объявление %s с названием %q не найдено за %d попыток",
			created.ID, title, searchAttempts)
	})

	s.Run("Проверяем загруженную фотографию", func() {
		require.Len(s.T(), created.Photos, 1, "При создании должна вернуться одна фотография")
		require.NotEmpty(s.T(), created.Photos[0].ID, "ID фотографии пустой")

		photos := advertisements.GetAdvertisementPhotos(s.T(), token, created.ID)
		require.Len(s.T(), photos, 1, "У объявления должна быть одна фотография")
		require.Equal(s.T(), created.Photos[0].ID, photos[0].ID, "ID фотографии не совпадает")
		require.NotEmpty(s.T(), photos[0].URL, "URL фотографии пустой")
	})
}
