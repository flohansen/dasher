package routes

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flohansen/dasher-server/internal/mocks"
	"github.com/flohansen/dasher-server/internal/model"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	featureStore := mocks.NewMockFeatureStore(ctrl)
	routes := New(featureStore)

	t.Run("GET /api/v1/features", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR", func(t *testing.T) {
			// given
			featureStore.EXPECT().
				GetAll().
				Return(nil, errors.New("some error"))

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/features", nil)
			routes.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			featureStore.EXPECT().
				GetAll().
				Return([]model.FeatureData{}, nil)

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/features", nil)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			b, _ := io.ReadAll(res.Body)
			assert.Equal(t, http.StatusOK, res.StatusCode)
			assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
			assert.JSONEq(t, `[]`, string(b))
		})
	})
}
