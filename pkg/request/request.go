package request

import (
	"errors"
	"log"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

var (
	Validate *validator.Validate
)

func init() {
	Validate = validator.New()
}

type BaseFindRequest struct {
	Limit int `form:"limit" validate:"max=500" default:"10"` // 获取多少条
	Skip  int `form:"skip"`                                  // 跳过多少条
}

func Assign(c *gin.Context, req interface{}) {
	if c.Request.Method == "GET" {
		if err := c.ShouldBindQuery(req); err != nil {
			log.Fatalf("ShouldBindQuery: %s", err.Error())
			panic(err)
		}
	} else {
		if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil && err.Error() != "EOF" {
			log.Fatalf("ShouldBindBodyWith: %s", err.Error())
			panic(err)
		}
	}
	// Set default values for fields with 'default' tag
	v := reflect.ValueOf(req).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		defaultValue := fieldType.Tag.Get("default")
		if defaultValue != "" && isEmptyValue(field) {
			setDefaultValue(field, defaultValue)
		}
	}

	if err := Validate.Struct(req); err != nil {
		log.Fatalf("Validate: %s", err.Error())
		panic(err)
	}
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}

func setDefaultValue(v reflect.Value, defaultValue string) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(defaultValue)
	case reflect.Bool:
		val, _ := strconv.ParseBool(defaultValue)
		v.SetBool(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(defaultValue, 10, 64)
		v.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		val, _ := strconv.ParseUint(defaultValue, 10, 64)
		v.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, _ := strconv.ParseFloat(defaultValue, 64)
		v.SetFloat(val)
	}
}

func CopyAndValidate(to interface{}, from interface{}) (err error) {
	if err = copier.CopyWithOption(to, from, copier.Option{IgnoreEmpty: true}); err != nil {
		return errors.New(err.Error())
	}
	if err := Validate.Struct(to); err != nil {
		return errors.New(err.Error())
	}
	return
}
