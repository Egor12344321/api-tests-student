package advertisementCreate

import (
	"net/http"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/managers/advertisements"
)

// негативные сценарии создания объявления с отсутствующим или некорректным токеном, но с валидными данными для объявления, чтобы изолированно проверить ошибку авторизации
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	const (
		unauthorizedError = "unauthorized"
		incorrectToken    = "incorrect_token"
	)

	cases := []struct {
		name    string
		token   string
		message string
	}{
		{"Без токена (негативный, заголовок Authorization отсутствует)", "", "Authorization header required"},
		{"С неверным токеном (негативный, токен недействителен)", incorrectToken, "Invalid or expired token"},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			fields := newAdvertisementFields()

			s.T().Log("Отправляем корректную форму и ожидаем HTTP 401")
			errorResponse := advertisements.CreateAdvertisementExpectError(
				s.T(), tc.token, fields, []string{photoPath}, http.StatusUnauthorized,
			)

			assertions.CheckErrorResponse(s.T(), errorResponse, unauthorizedError, tc.message)
		})
	}
}
