package entity

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

var (
	validate *validator.Validate
	initOne  = &sync.Once{}
)

func init() {
	initOne.Do(func() {
		validate = validator.New()

		// to get json tag using .Field() method
		validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})
}
