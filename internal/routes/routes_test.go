package routes

import (
	"bytes"
	"database/sql"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flohansen/dasher-server/internal/routes/mocks"
	"github.com/flohansen/dasher-server/internal/sqlc"
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
				GetAll(gomock.Any()).
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
				GetAll(gomock.Any()).
				Return(nil, nil)

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

	t.Run("POST /api/v1/features", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERROR if decoding request failed", func(t *testing.T) {
			// given
			body := bytes.NewReader([]byte(`{"feature`))

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/features", body)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		})

		t.Run("should return 500 INTERNAL SERVER ERROR if upserting feature failed", func(t *testing.T) {
			// given
			body := bytes.NewReader([]byte(`{"featureId":"TOGGLE_ID"}`))

			featureStore.EXPECT().
				Upsert(gomock.Any(), sqlc.Feature{
					FeatureID:   "TOGGLE_ID",
					Description: sql.NullString{String: "", Valid: false},
					Enabled:     sql.NullBool{Bool: false, Valid: false},
				}).
				Return(errors.New("some error"))

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/features", body)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			body := bytes.NewReader([]byte(`{"featureId":"TOGGLE_ID"}`))

			featureStore.EXPECT().
				Upsert(gomock.Any(), sqlc.Feature{
					FeatureID:   "TOGGLE_ID",
					Description: sql.NullString{String: "", Valid: false},
					Enabled:     sql.NullBool{Bool: false, Valid: false},
				}).
				Return(nil)

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/features", body)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	})

	t.Run("DELETE /api/v1/features", func(t *testing.T) {
		t.Run("should return 500 INTERNAL SERVER ERRROR if delete failed", func(t *testing.T) {
			// given
			featureStore.EXPECT().
				Delete(gomock.Any(), "TOGGLE_ID").
				Return(errors.New("some error"))

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/features/TOGGLE_ID", nil)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		})

		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			featureStore.EXPECT().
				Delete(gomock.Any(), "TOGGLE_ID").
				Return(nil)

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodDelete, "/api/v1/features/TOGGLE_ID", nil)
			routes.ServeHTTP(w, r)

			// then
			res := w.Result()
			assert.Equal(t, http.StatusOK, res.StatusCode)
		})
	})
}
