package customvalidation

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/leebenson/conform"
)

type ErrorWithTranslation struct {
	Err        error          `json:"err"`
	Translator *ut.Translator `json:"translator"`
}

// credits to this --> https://dev.to/koddr/how-to-make-clear-pretty-error-messages-from-the-go-backend-to-your-frontend-21b2

// NewValidator func for create a new validator for struct fields.
func NewValidator() *validator.Validate {
	return validator.New()
}

// ValidatorErrors func for show validation errors for each invalid fields.
func ValidatorErrors(err error) map[string]string {

	errFields := map[string]string{}

	for _, err := range err.(validator.ValidationErrors) {
		structName := strings.Split(err.Namespace(), ".")[0]
		errFields[err.Field()] = fmt.Sprintf(
			"failed '%s' tag check (value '%s' is not valid for %s struct)",
			err.Tag(), err.Value(), structName,
		)
	}
	return errFields
}

// run the sent payload through 2 functions
// 1. the BodyParser will check if the sent json object is correct
// 2 the NewValidator will check if the validate fields for the passed down struct are valid.
func ValidatePayload(c *gin.Context, payload interface{}) (err *ErrorWithTranslation) {

	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")

	validate := NewValidator()

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} should not be empty!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	if err := c.ShouldBindJSON(payload); err != nil {
		return &ErrorWithTranslation{
			Err:        err,
			Translator: &trans,
		}
	}

	if err := conform.Strings(payload); err != nil {
		return &ErrorWithTranslation{
			Err:        err,
			Translator: &trans,
		}
	}

	if err := validate.Struct(payload); err != nil {
		return &ErrorWithTranslation{
			Err:        err,
			Translator: &trans,
		}
	}
	return nil
}

func ValidateStruct(payload interface{}) (err *ErrorWithTranslation) {
	en := en.New()
	uni := ut.New(en, en)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("en")

	validate := NewValidator()

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} should not be empty!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	if err := conform.Strings(payload); err != nil {
		return &ErrorWithTranslation{
			Err:        err,
			Translator: &trans,
		}
	}

	if err := validate.Struct(payload); err != nil {
		return &ErrorWithTranslation{
			Err:        err,
			Translator: &trans,
		}
	}
	return nil
}
