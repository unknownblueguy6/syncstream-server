package valid8r

import "github.com/go-playground/validator/v10"

type ValidationErrors = validator.ValidationErrors

var Validator = validator.New(validator.WithRequiredStructEnabled())
