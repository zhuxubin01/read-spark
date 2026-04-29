package integration

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/readspark/backend/internal/handler"
	"github.com/readspark/backend/internal/service"
)

func TestAPISmoke_HealthMetricsAndDictionary(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[{"word":"hello","meanings":[{"definitions":[{"definition":"greeting"}]}]}]`))
	}))
	defer upstream.Close()

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	d := handler.NewDictionaryHandler(service.NewDictionaryServiceWithClient(upstream.URL, upstream.Client()))
	api := r.Group("/api/v1")
	api.GET("/dictionary/:word", d.Lookup)

	cases := []struct {
		path string
		code int
	}{
		{path: "/health", code: http.StatusOK},
		{path: "/metrics", code: http.StatusOK},
		{path: "/api/v1/dictionary/hello", code: http.StatusOK},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.code {
			t.Fatalf("path %s expected %d got %d", tc.path, tc.code, w.Code)
		}
	}
}
