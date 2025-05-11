package utils

import (
	"errors"
	"reflect"
)

// CopyNonNilFields copia los campos no nulos desde un struct origen (usualmente un DTO con punteros)
// hacia un destino (una entidad), campo por campo, usando reflection.
// Si un campo es nil en el origen y es de tipo string en el destino, se asigna "no information" como default.
func CopyNonNilFields(source interface{}, dest interface{}) error {
	srcVal := reflect.ValueOf(source)
	dstVal := reflect.ValueOf(dest)

	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return errors.New("parameters must be pointers to structs")
	}

	srcElem := srcVal.Elem()
	dstElem := dstVal.Elem()

	if srcElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Struct {
		return errors.New("parameters must be pointers to structs")
	}

	srcType := srcElem.Type()

	for i := 0; i < srcElem.NumField(); i++ {
		srcField := srcElem.Field(i)
		srcFieldType := srcType.Field(i)
		dstField := dstElem.FieldByName(srcFieldType.Name)

		if !dstField.IsValid() || !dstField.CanSet() {
			continue
		}

		if !srcField.IsNil() {
			// srcField es un puntero, tomamos su valor para setearlo
			dstField.Set(srcField.Elem())
		} else if dstField.Kind() == reflect.String {
			// Si el campo es string y está ausente, asignamos default
			dstField.SetString("no information")
		}
	}

	return nil
}