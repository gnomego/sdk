package validation

import "github.com/go-playground/validator/v10"

func IsSlug(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	return IsSlugValue(v)
}

func IsSlugValue(v string) bool {
	if len(v) == 0 {
		return true
	}

	for i, r := range v {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
			return false
		}

		if r == '-' && (i == 0 || i == len(v)-1) {
			return false
		}
	}

	return true
}
