package http

import (
	"strings"

	validatorV10 "github.com/go-playground/validator/v10"
)

type (
	ValidationErrorResponse struct {
		Errors []ValidationErrorItem `json:"errors"`
	}

	ValidationErrorItem struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}
)

func (v ValidationErrorResponse) String() string {
	var sb strings.Builder

	for i, err := range v.Errors {
		if i > 0 {
			sb.WriteString("; ")
		}
		sb.WriteString(err.Field)
		sb.WriteString(": ")
		sb.WriteString(err.Message)
	}
	return sb.String()
}

func ConvertValidationErrors(err error) ValidationErrorResponse {
	var response ValidationErrorResponse

	if errs, ok := err.(validatorV10.ValidationErrors); ok {
		for _, e := range errs {
			response.Errors = append(response.Errors, ValidationErrorItem{
				Field:   toSnakeCase(e.Namespace()),
				Message: messageForTag(e.Tag(), e.Param()),
			})
		}
	}

	return response
}

func messageForTag(tag, param string) string {
	switch tag {
	case "required":
		return "field is required"
	case "min", "gt", "gte":
		return "value is too short, min=" + param
	case "max", "lt", "lte":
		return "value is too long, max=" + param
	case "email":
		return "invalid email format"
	default:
		return "invalid value"
	}
}

func toSnakeCase(ns string) string {
	parts := strings.SplitN(ns, ".", 2)
	if len(parts) == 2 {
		ns = parts[1]
	}
	return strings.ToLower(ns)
}
