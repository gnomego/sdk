package validation

import (
	"reflect"
	"sync"

	"github.com/gnomego/apps/gs/xgin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type GsValidator struct {
	*validator.Validate
	uni  *ut.UniversalTranslator
	once sync.Once
}

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func NewGsValidator() *GsValidator {
	eng := en.New()
	translater := ut.New(eng, eng)
	v := &GsValidator{
		Validate: validator.New(),
		uni:      translater,
		once:     sync.Once{},
	}

	v.RegisterGsValidation()

	return v
}

func (v *GsValidator) RegisterGsValidation() *GsValidator {
	v.RegisterValidation("slug", IsSlug)
	v.RegisterValidation("domain", IsDomain)
	v.AddTranslation("en", "slug", "{0} must be a valid slug with only lowercase letters, numbers, and hyphens.")
	v.AddTranslation("en", "domain", "{0} must be a valid domain with only lowercase letters, numbers, hyphens, and periods.")
	return v
}

func (v *GsValidator) TranslateStruct(obj interface{}, err error) *xgin.ErrorInfo {
	if kindOfData(obj) != reflect.Struct {
		return nil
	}

	structName := reflect.TypeOf(obj).Name()
	return v.TranslateStructWithName(obj, structName, err)
}

func (v *GsValidator) TranslateStructWithName(obj interface{}, target string, err error) *xgin.ErrorInfo {
	if err == nil {
		return nil
	}

	if kindOfData(obj) != reflect.Struct {
		return nil
	}

	root := xgin.NewErrorInfo("validation", "One or more validation errors occurred")
	root.SetTarget(target)

	translator, _ := v.uni.FindTranslator("en")
	errs := err.(validator.ValidationErrors)
	details := make([]xgin.ErrorInfo, 0)
	for _, e := range errs {
		message := e.Translate(translator)
		field := e.Field()

		next := xgin.NewErrorInfo(e.Tag(), message).SetTarget(field)
		details = append(details, *next)
	}

	root.Details = details
	return root
}

func New() *DefaultValidator {
	v := &DefaultValidator{}
	v.validate = validator.New()
	v.validate.RegisterValidation("slug", IsSlug)
	return v
}

func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func (v *GsValidator) AddEngTranslation(tag string, message string) {
	v.AddTranslation("en", tag, message)
}

func (v *GsValidator) OverrideEngTranslation(tag string, message string) {
	v.OverrideTranslation("en", tag, message)
}

func (v *GsValidator) OverrideTranslation(lang string, tag string, message string) {
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

func (v *GsValidator) AddTranslation(lang string, tag string, message string) {
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
