package validation

import (
	"reflect"
	"sync"

	"github.com/gnomego/sdk/errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var globalValidator *Validator

type Validator struct {
	*validator.Validate
	uni           *ut.UniversalTranslator
	once          sync.Once
	defaultsAdded bool
}

func New() *Validator {
	eng := en.New()
	translater := ut.New(eng, eng)
	v := &Validator{
		Validate: validator.New(validator.WithRequiredStructEnabled()),
		uni:      translater,
		once:     sync.Once{},
	}

	return v
}

func Default() *Validator {
	if globalValidator == nil {
		globalValidator = New()
	}

	return globalValidator
}

func (v *Validator) RegisterExtraValidators() *Validator {
	if !v.defaultsAdded {
		v.defaultsAdded = true
		english, ok := v.uni.GetTranslator("en")
		if !ok {
			english = v.uni.GetFallback()
		}

		en_translations.RegisterDefaultTranslations(v.Validate, english)
		v.RegisterValidation("slug", IsSlug)
		v.RegisterValidation("domain", IsDomain)
		v.AddTranslation("en", "slug", "{0} must be a valid slug with only lowercase letters, numbers, and hyphens.")
		v.AddTranslation("en", "domain", "{0} must be a valid domain with only lowercase letters, numbers, hyphens, and periods.")
	}

	return v
}

func (v *Validator) TranslateStruct(obj interface{}, err error) *errors.SystemError {
	if kindOfData(obj) != reflect.Struct {
		return nil
	}

	structName := reflect.TypeOf(obj).Name()
	return v.TranslateStructWithName(obj, structName, err)
}

func (v *Validator) TranslateStructWithName(obj interface{}, target string, err error) *errors.SystemError {
	if err == nil {
		return nil
	}

	if kindOfData(obj) != reflect.Struct {
		return nil
	}

	root := errors.NewWithCodef("validation_error", "One or more errors found for %s", target)
	root.SetTarget(target)

	translator, _ := v.uni.FindTranslator("en")
	errs := err.(validator.ValidationErrors)
	for _, e := range errs {
		message := e.Translate(translator)
		field := e.Field()

		root.AddDetail(errors.NewWithCodef(e.Tag(), message).SetTarget(field))
	}

	return root
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *Validator) AddEnglishTranslation(tag string, message string) {
	v.AddTranslation("en", tag, message)
}

func (v *Validator) OverrideEnglishTranslation(tag string, message string) {
	v.OverrideTranslation("en", tag, message)
}

func (v *Validator) OverrideTranslation(lang string, tag string, message string) {
	translator, ok := v.uni.GetTranslator(lang)
	if !ok {
		translator = v.uni.GetFallback()
	}

	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, message, true)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = v.RegisterTranslation(tag, translator, registerFn, transFn)
}

func (v *Validator) init() {
	if v.uni == nil {
		english := en.New()
		v.uni = ut.New(english, english)
	}
}

func (v *Validator) AddTranslation(lang string, tag string, message string) {
	v.init()

	translator, ok := v.uni.GetTranslator(lang)
	if !ok {
		translator = v.uni.GetFallback()
	}

	registerFn := func(ut ut.Translator) error {
		return ut.Add(tag, message, false)
	}

	transFn := func(ut ut.Translator, fe validator.FieldError) string {
		param := fe.Param()
		tag := fe.Tag()

		t, err := ut.T(tag, fe.Field(), param)
		if err != nil {
			return fe.(error).Error()
		}
		return t
	}

	_ = v.RegisterTranslation(tag, translator, registerFn, transFn)
}
