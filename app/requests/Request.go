package requests

import (
	"reflect"
	"github.com/go-playground/validator/v10"
)

type RequestHelper struct{}

func NewRequestHelper() *RequestHelper {
	return &RequestHelper{}
}

// GetJSONTag retorna a tag JSON para um campo específico
func (rh *RequestHelper) GetJSONTag(field string, structs ...interface{}) (string, bool) {
	for _, s := range structs {
		rv := reflect.ValueOf(s)
		for i := 0; i < rv.NumField(); i++ {
			if rv.Type().Field(i).Name == field {
				tag := rv.Type().Field(i).Tag.Get("json")
				if tag != "" {
					return tag, true
				}
			}
		}
	}
	return "", false
}

//mensagems de erro customizadas
func GetErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "Este campo é obrigatório"
	case "email":
		return "Deve ser um e-mail válido"
	case "min":
		return "Deve ter pelo menos " + err.Param() + " caracteres"
	case "eqfield":
		return "Deve ser igual a " + err.Param()
	default:
		return "Entrada inválida"
	}
}
