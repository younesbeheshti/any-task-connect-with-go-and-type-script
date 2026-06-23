package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	phoneRegex        = regexp.MustCompile(`^09\d{9}$`)
	nationalCodeRegex = regexp.MustCompile(`^\d{10}$`)
)

// Validator wraps go-playground/validator with custom rules.
type Validator struct {
	validate *validator.Validate
}

// New creates a Validator with custom validations registered.
func New() *Validator {
	v := validator.New()
	_ = v.RegisterValidation("phone", validatePhone)
	_ = v.RegisterValidation("national_code", validateNationalCode)
	_ = v.RegisterValidation("password_strength", validatePasswordStrength)
	return &Validator{validate: v}
}

// ValidateStruct validates a struct using validation tags.
func (v *Validator) ValidateStruct(s any) error {
	if err := v.validate.Struct(s); err != nil {
		return formatValidationError(err)
	}
	return nil
}

// Engine returns the underlying validator engine.
func (v *Validator) Engine() *validator.Validate {
	return v.validate
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func validateNationalCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	if !nationalCodeRegex.MatchString(code) {
		return false
	}
	return isValidIranianNationalCode(code)
}

func validatePasswordStrength(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) >= 8
}

func isValidIranianNationalCode(code string) bool {
	if code == "0000000000" {
		return false
	}
	check := int(code[9] - '0')
	sum := 0
	for i := 0; i < 9; i++ {
		sum += int(code[i]-'0') * (10 - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		return check == remainder
	}
	return check == 11-remainder
}

func formatValidationError(err error) error {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	fields := make(map[string]string, len(validationErrs))
	for _, fe := range validationErrs {
		field := strings.ToLower(fe.Field())
		fields[field] = validationMessage(fe)
	}
	return ValidationError{Fields: fields}
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "phone":
		return "must be a valid Iranian mobile number (09XXXXXXXXX)"
	case "national_code":
		return "must be a valid national code"
	case "password_strength":
		return "رمز عبور باید حداقل ۸ کاراکتر باشد"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return fmt.Sprintf("failed on '%s' validation", fe.Tag())
	}
}

// ValidationError represents field-level validation failures.
type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string {
	return "validation failed"
}
