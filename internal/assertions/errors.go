package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"
)

// CheckErrorResponse проверяет код и сообщение ошибки API.
func CheckErrorResponse(t *testing.T, errorResponse models.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()

	require.Equal(t, expectedError, errorResponse.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, errorResponse.Message, "Некорректное значение message параметра")
}
