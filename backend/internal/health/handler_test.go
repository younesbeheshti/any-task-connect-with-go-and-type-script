package health_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/internal/health"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/testutil"
)

type alwaysUp struct{}

func (alwaysUp) HealthCheck(context.Context) error { return nil }

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := health.NewHandler(
		testutil.MockDB(t),
		alwaysUp{},
		alwaysUp{},
	)

	r := gin.New()
	handler.RegisterRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, "UP", body["status"])

	services, ok := body["services"].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "UP", services["database"])
	assert.Equal(t, "UP", services["redis"])
	assert.Equal(t, "UP", services["rabbitmq"])
}
