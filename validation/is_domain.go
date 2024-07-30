package validation

import "github.com/go-playground/validator/v10"

func IsDomain(fl validator.FieldLevel) bool {
	v := fl.Field().String()

	if len(v) == 0 {
		return true
	}

	// check domain length
	if len(v) > 253 {
		return false
	}

	for i, r := range v {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '.' {
			return false
		}

		if r == '-' && (i == 0 || i == len(v)-1) {
			return false
		}

		if r == '.' && (i == 0 || i == len(v)-1) {
			return false
		}
	}

	return true
}
