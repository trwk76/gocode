package code

import "github.com/iancoleman/strcase"

type Casing string

const (
	CasingNone   Casing = ""
	CasingPascal Casing = "pascal"
	CasingCamel  Casing = "camel"
	CasingSnake  Casing = "snake"
)

func (c Casing) To(name string) string {
	switch c {
	case CasingPascal:
		return strcase.ToCamel(name)
	case CasingCamel:
		return strcase.ToLowerCamel(name)
	case CasingSnake:
		return strcase.ToSnake(name)
	}

	return name
}
