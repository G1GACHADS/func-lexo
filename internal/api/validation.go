package api

import (
	"fmt"
	"reflect"
)

func RequiredFields(fields map[string]any) error {
	for k, v := range fields {
		if err := RequiredField(k, v); err != nil {
			return err
		}
	}

	return nil
}

func RequiredField(fieldName string, v any) error {
	if reflect.ValueOf(v).IsZero() {
		return fmt.Errorf("%s is required", fieldName)
	}

	return nil
}
