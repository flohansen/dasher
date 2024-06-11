package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	t.Run("GET /api/v1/features", func(t *testing.T) {
		t.Run("should return 200 OK", func(t *testing.T) {
			// given
			routes := New()

			// when
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/features", nil)
			routes.ServeHTTP(w, r)

			// then
			assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		})
	})
}
