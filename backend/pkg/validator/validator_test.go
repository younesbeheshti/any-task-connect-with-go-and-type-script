package validator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/younesbeheshti/any-task-connect/backend/pkg/validator"
)

func TestPasswordStrength(t *testing.T) {
	v := validator.New()

	err := v.ValidateStruct(struct {
		Password string `validate:"password_strength"`
	}{Password: "weak"})
	require.Error(t, err)

	err = v.ValidateStruct(struct {
		Password string `validate:"password_strength"`
	}{Password: "Strong1!"})
	assert.NoError(t, err)
}

func TestNationalCode(t *testing.T) {
	v := validator.New()
	err := v.ValidateStruct(struct {
		Code string `validate:"national_code"`
	}{Code: "0000000000"})
	require.Error(t, err)
}
