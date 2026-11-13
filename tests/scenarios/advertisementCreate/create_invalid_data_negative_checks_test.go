package advertisementCreate

import (
	"net/http"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

// параметризированный тест с негативными проверками отсутствия обязательных полей, с граничными значениями и классами эквиваленности без проверки поиска по названию и допуступности фото
func (s *TestSuite) TestCreateAdvertisementInvalidFields() {
	const (
		invalidRequest        = "invalid_request"
		titleRequired         = "title is required"
		descriptionRequired   = "description is required"
		descriptionRangeError = "description must be between 1 and 500 characters"
		priceRangeError       = "price must be between 0 and 1000000 inclusive"
		quantityRangeError    = "quantity must be between 1 and 100 inclusive"
		photoRequired         = "at least 1 photo is required"
		photoLimitError       = "maximum 3 photos allowed"
	)

	cases := []struct {
		name       string
		field      string
		value      string
		omitField  bool
		photoCount int
		message    string
	}{
		{"Название отсутствует (негативный, обязательность поля)", "title", "", true, 1, titleRequired},
		{"Название пустое (негативный, граничное значение, 0 при минимуме 1 символ)", "title", "", false, 1, titleRequired},
		{"Название из 51 символа (негативный, граничное значение, максимум 50)", "title", strings.Repeat("A", 51), false, 1, "title must not exceed 50 characters provided 51"},
		{"Название из 60 символов (негативный, класс эквивалентности, больше 50 символов)", "title", strings.Repeat("A", 60), false, 1, "title must not exceed 50 characters provided 60"},

		{"Описание отсутствует (негативный, обязательность поля)", "description", "", true, 1, descriptionRequired},
		{"Описание пустое (негативный, граничное значение, 0 при минимуме 1 символ)", "description", "", false, 1, descriptionRequired},
		{"Описание из 501 символа (негативный, граничное значение, максимум 500)", "description", strings.Repeat("A", 501), false, 1, descriptionRangeError},
		{"Описание из 600 символов (негативный, класс эквивалентности, больше 500 символов)", "description", strings.Repeat("A", 600), false, 1, descriptionRangeError},

		{"Цена −1 (негативный, граничное значение с шагом 1, минимум 0)", "price", "-1", false, 1, priceRangeError},
		{"Цена −100 (негативный, класс эквивалентности, цена ниже 0)", "price", "-100", false, 1, priceRangeError},
		{"Цена 1000001 (негативный, граничное значение с шагом 1, максимум 1000000)", "price", "1000001", false, 1, priceRangeError},
		{"Цена 1500000 (негативный, класс эквивалентности, цена выше 1000000)", "price", "1500000", false, 1, priceRangeError},

		{"Количество 0 (негативный, граничное значение, минимум 1)", "quantity", "0", false, 1, quantityRangeError},
		{"Количество −10 (негативный, класс эквивалентности, количество меньше 1)", "quantity", "-10", false, 1, quantityRangeError},
		{"Количество 101 (негативный, граничное значение, максимум 100)", "quantity", "101", false, 1, quantityRangeError},
		{"Количество 150 (негативный, класс эквивалентности, количество больше 100)", "quantity", "150", false, 1, quantityRangeError},

		{"0 фотографий (негативный, граничное значение, минимум 1 файл)", "", "", false, 0, photoRequired},
		{"4 фотографии (негативный, граничное значение, максимум 3 файла)", "", "", false, 4, photoLimitError},
		{"5 фотографий (негативный, класс эквивалентности, больше 3 файлов)", "", "", false, 5, photoLimitError},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			fields := newAdvertisementFields()
			if tc.omitField {
				delete(fields, tc.field)
			} else if tc.field != "" {
				fields[tc.field] = tc.value
			}

			s.T().Log("Отправляем запрос и ожидаем HTTP 400")
			errorResponse := advertisements.CreateAdvertisementExpectError(
				s.T(), s.loginData.Token, fields, advertisementPhotoPaths(tc.photoCount), http.StatusBadRequest,
			)

			require.Equal(s.T(), invalidRequest, errorResponse.Error, "Код ошибки не совпадает")
			require.Equal(s.T(), tc.message, errorResponse.Message, "Сообщение об ошибке не совпадает")
		})
	}
}

// advertisementPhotoPaths возвращает пути к тестовому PNG для нужного количества фотографий.
func advertisementPhotoPaths(count int) []string {
	photoPaths := make([]string, count)
	for i := range photoPaths {
		photoPaths[i] = photoPath
	}
	return photoPaths
}
