package validate

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
	enTranslations "gopkg.in/go-playground/validator.v9/translations/en"
	"reflect"
)

// field validator for post method
var bindingValidator *validator.Validate
// field validator for patch or put method
var updateValidator *validator.Validate

// validate err translator
var BindingTrans ut.Translator
var UpdateTrans ut.Translator

func init() {
	bindingValidator, _ = binding.Validator.Engine().(*validator.Validate)
	en0 := en.New()
	uni0 := ut.New(en0, en0)
	// register err message English
	BindingTrans, _ = uni0.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(bindingValidator, BindingTrans)
	bindingValidator.RegisterTagNameFunc(registerTagName)

	updateValidator = validator.New()
	updateValidator.SetTagName("update")
	en1 := en.New()
	uni1 := ut.New(en1, en1)
	// register err message English
	UpdateTrans, _ = uni1.GetTranslator("en")
	enTranslations.RegisterDefaultTranslations(updateValidator, UpdateTrans)
	updateValidator.RegisterTagNameFunc(registerTagName)
	updateValidator.RegisterValidation("fixed", fixed)
	updateValidator.RegisterTranslation("fixed", UpdateTrans, fixedTranslation, fixedTranslationAdding)
}

func registerTagName(field reflect.StructField) string {
	name := field.Tag.Get("json")
	if name == "-" {
		return ""
	}
	return name
}

//this field must be read only, can not be update
func fixed(f1 validator.FieldLevel) bool {
	return false
}

// fixed tag custom error message
func fixedTranslation(ut ut.Translator) error {
	return ut.Add("fixed", "{0} can not be updated. it's fixed.", true)
}

// register fixed tag custom error message
func fixedTranslationAdding(ut ut.Translator, fe validator.FieldError) string {
	t, _ := ut.T("fixed", fe.Field())
	return t
}

//obj must be struct only, other type is not allowed
func StructForUpdate(obj interface{}, structFieldNames map[string]bool) error {
	if len(structFieldNames) == 0 {
		return errors.New("no attribute is provided")
	}
	return updateValidator.Struct(obj)
}
