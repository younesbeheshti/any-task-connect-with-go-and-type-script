package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/configs"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/jwt"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
	svc := jwt.NewService(configs.JWTConfig{
		Secret: "test-secret-key-minimum-32-characters-long",
		AccessTTL: 900, RefreshTTL: 604800,
	})

	token, err := svc.GenerateAccessToken("user-1", common.RoleRequester)
	require.NoError(t, err)

	claims, err := svc.ValidateAccessToken(token)
	require.NoError(t, err)
	assert.Equal(t, "user-1", claims.UserID)
	assert.Equal(t, common.RoleRequester, claims.Role)
}

func TestRefreshTokenRotation(t *testing.T) {
	svc := jwt.NewService(configs.JWTConfig{
		Secret: "test-secret-key-minimum-32-characters-long",
		AccessTTL: 900, RefreshTTL: 604800,
	})

	refresh, err := svc.GenerateRefreshToken("user-1", common.RoleAgent)
	require.NoError(t, err)

	access, newRefresh, err := svc.RefreshToken(refresh)
	require.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, newRefresh)
}
