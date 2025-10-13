package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := &CustomValidator{
		validator: validator.New(),
	}

	v.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		jsonName := getJSONFieldName(field)
		if jsonName == "-" {
			return ""
		}
		return jsonName
	})

	return v
}

func (v *CustomValidator) Validate(data interface{}) error {
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return FieldErrors{{
			Field:   "data",
			Value:   data,
			Message: "Data must be a pointer",
		}}
	}

	errors := validateAndExtractErrors(v.validator, data)
	if len(errors) == 0 {
		return nil
	}

	return errors
}

func validateAndExtractErrors(validation *validator.Validate, data interface{}) FieldErrors {
	var errorDetails FieldErrors

	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		message := "Parameter data must be a pointer"

		return FieldErrors{{
			Field:   "data",
			Value:   data,
			Message: message,
		}}
	}

	err := validation.Struct(data)
	if err == nil {
		return errorDetails
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		message := "Unknown validation error"

		return FieldErrors{{
			Field:   "unknown",
			Value:   nil,
			Message: message,
		}}
	}

	dataValue := reflect.ValueOf(data)
	dataType := reflect.TypeOf(data).Elem()

	for _, validationErr := range validationErrs {
		fieldError := processValidationError(validationErr, dataValue, dataType)
		errorDetails = append(errorDetails, fieldError)
	}

	return errorDetails
}

func processValidationError(
	validationErr validator.FieldError,
	dataValue reflect.Value,
	dataType reflect.Type,
) FieldError {
	fieldName := validationErr.StructField()
	jsonName, fieldValue := extractFieldInfo(fieldName, dataValue, dataType)
	errorMessage := translateValidationTag(validationErr.Tag(), validationErr.Param())

	return FieldError{
		Field:   jsonName,
		Value:   fieldValue,
		Message: errorMessage,
	}
}

func extractFieldInfo(fieldName string, dataValue reflect.Value, dataType reflect.Type) (string, interface{}) {
	field, found := dataType.FieldByName(fieldName)

	if !found {
		return fieldName, "<unknown>"
	}

	jsonName := getJSONFieldName(field)
	if jsonName == "" {
		jsonName = fieldName
	}

	fieldValue := extractFieldValue(fieldName, dataValue)

	return jsonName, fieldValue
}

func extractFieldValue(fieldName string, dataValue reflect.Value) interface{} {
	if dataValue.Kind() != reflect.Ptr || dataValue.IsNil() {
		return nil
	}

	elemValue := dataValue.Elem()
	if elemValue.Kind() != reflect.Struct {
		return nil
	}

	fieldValue := elemValue.FieldByName(fieldName)
	if !fieldValue.IsValid() {
		return nil
	}

	if !fieldValue.CanInterface() {
		return getDefaultValue(fieldValue.Kind())
	}

	return fieldValue.Interface()
}

func getDefaultValue(kind reflect.Kind) interface{} {
	switch kind {
	case reflect.String:
		return ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(0)
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.Bool:
		return false
	case reflect.Slice, reflect.Map:
		return "<empty>"
	default:
		return "<nil>"
	}
}

func getJSONFieldName(field reflect.StructField) string {
	jsonTagValue := field.Tag.Get("json")
	if jsonTagValue == "" {
		return field.Name
	}

	parts := strings.SplitN(jsonTagValue, ",", 2)
	if parts[0] == "" {
		return field.Name
	}

	return parts[0]
}

func translateValidationTag(tag string, param string) string {
	var translationKey string
	switch tag {
	case "required":
		translationKey = "validation_required"
	case "email":
		translationKey = "validation_email"
	case "min":
		translationKey = "validation_min"
	case "max":
		translationKey = "validation_max"
	case "len":
		translationKey = "validation_len"
	case "gt":
		translationKey = "validation_gt"
	case "lt":
		translationKey = "validation_lt"
	case "gte":
		translationKey = "validation_gte"
	case "lte":
		translationKey = "validation_lte"
	case "oneof":
		translationKey = "validation_oneof"
	case "alphanum":
		translationKey = "validation_alphanum"
	case "numeric":
		translationKey = "validation_numeric"
	case "alpha":
		translationKey = "validation_alpha"
	default:
		translationKey = "validation_generic"
	}

	message := translationKey

	if message == translationKey {
		switch tag {
		case "required":
			return "This field is required."
		case "email":
			return "Must be a valid email address."
		case "min":
			return fmt.Sprintf("Must be at least %s.", param)
		case "max":
			return fmt.Sprintf("Must not exceed %s.", param)
		case "len":
			return fmt.Sprintf("Must be exactly %s characters long.", param)
		case "gt":
			return fmt.Sprintf("Must be greater than %s.", param)
		case "lt":
			return fmt.Sprintf("Must be less than %s.", param)
		case "gte":
			return fmt.Sprintf("Must be greater than or equal to %s.", param)
		case "lte":
			return fmt.Sprintf("Must be less than or equal to %s.", param)
		case "oneof":
			return fmt.Sprintf("Must be one of the following: [%s].", param)
		case "alphanum":
			return "Must contain only alphanumeric characters."
		case "numeric":
			return "Must contain only numeric values."
		case "alpha":
			return "Must contain only alphabetic characters."
		default:
			return fmt.Sprintf("Validation failed for tag: %s.", tag)
		}
	}

	// Replace parameter placeholder if needed
	if strings.Contains(message, "{param}") {
		message = strings.Replace(message, "{param}", param, -1)
	}

	return message
}
