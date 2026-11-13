package advertisementCreate

import (
	"strconv"
	"strings"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements"
)

// параметризированный тест с позитивными проверками с граничными значениями и классами эквиваленности без проверки поиска по названию и допуступности фото
func (s *TestSuite) TestCreateAdvertisementValidFields() {
	cases := []struct {
		name       string
		field      string
		value      string
		omitField  bool
		photoCount int
	}{
		{"Название из 1 символа (позитивный, граничное значение, минимум 1 символ)", "title", "A", false, 1},
		{"Название из 25 символов (позитивный, класс эквивалентности, допустимая длина 1–50)", "title", strings.Repeat("A", 25), false, 1},
		{"Название из 49 символов (позитивный, граничное значение, на 1 меньше максимума 50)", "title", strings.Repeat("A", 49), false, 1},
		{"Название из 50 символов (позитивный, граничное значение, максимум 50)", "title", strings.Repeat("A", 50), false, 1},

		{"Описание из 1 символа (позитивный, граничное значение, минимум 1 символ)", "description", "A", false, 1},
		{"Описание из 250 символов (позитивный, класс эквивалентности, допустимая длина 1–500)", "description", strings.Repeat("A", 250), false, 1},
		{"Описание из 499 символов (позитивный, граничное значение, на 1 меньше максимума 500)", "description", strings.Repeat("A", 499), false, 1},
		{"Описание из 500 символов (позитивный, граничное значение, максимум 500)", "description", strings.Repeat("A", 500), false, 1},

		{"Цена 0 (позитивный, граничное значение, минимум 0)", "price", "0", false, 1},
		{"Цена 123.45 (позитивный, класс эквивалентности, допустимая дробная цена)", "price", "123.45", false, 1},
		{"Цена 999999 (позитивный, граничное значение с шагом 1, максимум 1000000)", "price", "999999", false, 1},
		{"Цена 1000000 (позитивный, граничное значение, максимум 1000000)", "price", "1000000", false, 1},
		{"Цена отсутствует (позитивный, необязательное поле)", "price", "", true, 1},

		{"Количество 1 (позитивный, граничное значение, минимум 1)", "quantity", "1", false, 1},
		{"Количество 45 (позитивный, класс эквивалентности, допустимое количество 1–100)", "quantity", "45", false, 1},
		{"Количество 99 (позитивный, граничное значение, на 1 меньше максимума 100)", "quantity", "99", false, 1},
		{"Количество 100 (позитивный, граничное значение, максимум 100)", "quantity", "100", false, 1},
		{"Количество отсутствует (позитивный, необязательное поле)", "quantity", "", true, 1},

		{"2 фотографии (позитивный, класс эквивалентности, допустимо от 1 до 3 файлов)", "", "", false, 2},
		{"3 фотографии (позитивный, граничное значение, максимум 3 файла)", "", "", false, 3},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			t := s.T()
			token := s.loginData.Token
			fields := newAdvertisementFields()
			if tc.omitField {
				delete(fields, tc.field)
			} else if tc.field != "" {
				fields[tc.field] = tc.value
			}

			t.Log("Создаём объявление с допустимыми данными")
			created := advertisements.CreateAdvertisement(t, token, fields, advertisementPhotoPaths(tc.photoCount))
			t.Cleanup(func() {
				t.Logf("Удаляем объявление %s", created.ID)
				advertisements.DeleteAdvertisement(t, token, created.ID)
			})

			t.Log("Проверяем сохранённые поля и количество фотографий")
			saved := advertisements.GetAdvertisementByID(t, created.ID)
			require.Equal(t, created.ID, saved.ID, "ID объявления не совпадает")

			savedFields := map[string]string{
				"title":       saved.Title,
				"description": saved.Description,
				"price":       strconv.FormatFloat(saved.Price, 'f', -1, 64),
				"quantity":    strconv.Itoa(saved.Quantity),
			}
			for field, value := range fields {
				require.Equalf(t, value, savedFields[field], "Поле %q не совпадает", field)
			}
			require.Len(t, created.Photos, tc.photoCount, "Количество фотографий не совпадает")
		})
	}
}
