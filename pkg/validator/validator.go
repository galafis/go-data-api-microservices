package validator

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	// Validator is the global validator instance
	validate *validator.Validate
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

// ValidationErrors represents a collection of validation errors
type ValidationErrors []ValidationError

// Error returns the error message
func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	for i, err := range ve {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Message)
	}
	return sb.String()
}

// init initializes the validator
func init() {
	validate = validator.New()

	// Register custom validation tags
	_ = validate.RegisterValidation("uuid", validateUUID)
	_ = validate.RegisterValidation("alpha_space", validateAlphaSpace)
	_ = validate.RegisterValidation("phone", validatePhone)
	_ = validate.RegisterValidation("password", validatePassword)
}

// Validate validates a struct
func Validate(i interface{}) error {
	err := validate.Struct(i)
	if err == nil {
		return nil
	}

	// Convert validation errors
	var validationErrors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		tag := err.Tag()
		value := fmt.Sprintf("%v", err.Value())
		message := getErrorMessage(field, tag, value, err.Param())

		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Tag:     tag,
			Value:   value,
			Message: message,
		})
	}

	return validationErrors
}

// ValidateVar validates a variable
func ValidateVar(field interface{}, tag string) error {
	err := validate.Var(field, tag)
	if err == nil {
		return nil
	}

	// Convert validation errors
	var validationErrors ValidationErrors
	for _, err := range err.(validator.ValidationErrors) {
		fieldName := reflect.TypeOf(field).Name()
		tag := err.Tag()
		value := fmt.Sprintf("%v", err.Value())
		message := getErrorMessage(fieldName, tag, value, err.Param())

		validationErrors = append(validationErrors, ValidationError{
			Field:   fieldName,
			Tag:     tag,
			Value:   value,
			Message: message,
		})
	}

	return validationErrors
}

// Custom validators

// validateUUID validates a UUID
func validateUUID(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	uuidRegex := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	return uuidRegex.MatchString(strings.ToLower(value))
}

// validateAlphaSpace validates a string containing only letters and spaces
func validateAlphaSpace(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	alphaSpaceRegex := regexp.MustCompile(`^[a-zA-Z\s]+$`)
	return alphaSpaceRegex.MatchString(value)
}

// validatePhone validates a phone number
func validatePhone(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	phoneRegex := regexp.MustCompile(`^\+?[0-9]{10,15}$`)
	return phoneRegex.MatchString(value)
}

// validatePassword validates a password
func validatePassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if value == "" {
		return true
	}

	// At least 8 characters, one uppercase, one lowercase, one number, one special character
	passwordRegex := regexp.MustCompile(`^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`)
	return passwordRegex.MatchString(value)
}

// getErrorMessage returns a human-readable error message for a validation error
func getErrorMessage(field, tag, value, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, param)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "alpha_space":
		return fmt.Sprintf("%s must contain only letters and spaces", field)
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", field)
	case "password":
		return fmt.Sprintf("%s must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number, and one special character", field)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	default:
		return fmt.Sprintf("%s failed validation for tag %s", field, tag)
	}
}

