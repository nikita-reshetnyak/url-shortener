package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSaveHandler(t *testing.T) {
	cases := []struct {
		name      string // Имя теста
		alias     string // Отправляемый alias
		url       string // Отправляемый URL
		respError string // Какую ошибку мы должны получить?
		mockError error  // Ошибку, которую вернёт мок
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaverMock := mocks.NewUrlSaver(t)
			if tc.respError == "" || tc.mockError != nil {
				urlSaverMock.On("SaveUrl", tc.url, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).Once()
			}
			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaverMock)
			payload := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias)
			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(payload)))
			require.NoError(t, err)

			resprec := httptest.NewRecorder()
			handler.ServeHTTP(resprec, req)

			require.Equal(t, resprec.Code, http.StatusOK)
			body := resprec.Body.String()

			var resp save.Response
			require.NoError(t, json.Unmarshal([]byte(body), &resp))
			require.Equal(t, tc.respError, resp.Error)
		})
	}

}
