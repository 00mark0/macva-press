package api

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// CustomValidator wraps the validator package and is used by Echo for request validation.
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator initializes the validator, registers any custom validations,
// and returns a new instance of CustomValidator.
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Example of registering a custom validation:
	// v.RegisterValidation("currency", validCurrency)
	// You can add more custom validations here as needed.
	// Register custom validations
	v.RegisterValidation("regex", validCategoryName)
	v.RegisterValidation("password", validPassword)
	v.RegisterValidation("username", validUsername)

	return &CustomValidator{validator: v}
}

// Validate implements the echo.Validator interface.
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return err
	}
	return nil
}

// If needed, you can add helper functions for additional validations here.
// For instance, if you want a custom "currency" validation:
//
// var validCurrency validator.Func = func(fl validator.FieldLevel) bool {
//     // Add your logic to check if the currency is supported, e.g.,
//     // return utils.IsSupportedCurrency(fl.Field().String())
//     return true // Simplified example
// }

// validCategoryName is a custom validation function for the "regex" tag
var validCategoryName validator.Func = func(fl validator.FieldLevel) bool {
	// Regex pattern: allow letters (both Latin and Cyrillic) and spaces
	re := regexp.MustCompile(`^[A-Za-zА-Яа-яČčĆćŽž ]+$`)
	return re.MatchString(fl.Field().String())
}

// validPassword enforces strong password rules
var validPassword validator.Func = func(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check length constraint
	if len(password) < 8 || len(password) > 64 {
		return false
	}

	var hasLower, hasUpper, hasDigit, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

// validUsername enforces the validation rules for the "username" tag.
var validUsername validator.Func = func(fl validator.FieldLevel) bool {
	username := fl.Field().String()

	// Check if the username length is between 3 and 20 characters
	if len(username) < 3 || len(username) > 20 {
		return false
	}

	// Check if the username starts with a letter
	if unicode.IsDigit(rune(username[0])) {
		return false
	}

	// Check each character to ensure it only contains valid characters (alphanumeric, _, or -)
	for _, char := range username {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' && char != '-' {
			return false
		}
	}

	// Username is valid if it passes all checks
	return true
}
