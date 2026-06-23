package common_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/younesbeheshti/any-task-connect/backend/internal/common"
)

func TestHasPermission(t *testing.T) {
	assert.True(t, common.HasPermission(common.RoleRequester, common.PermTaskCreate))
	assert.False(t, common.HasPermission(common.RoleAgent, common.PermTaskCreate))
	assert.True(t, common.HasPermission(common.RoleAdmin, common.PermAdminDashboard))
}

func TestRoleFromAPI(t *testing.T) {
	assert.Equal(t, common.RoleRequester, common.RoleFromAPI("requester"))
	assert.Equal(t, common.RoleAgent, common.RoleFromAPI("agent"))
}
