package db

import (
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func expressionByField(db *gorm.DB, query interface{}) *gorm.DB {
	// 通过反射获取查询条件
	// 获取字段的标签信息 column 数据库列名 operate 操作符
	queryType := reflect.TypeOf(query)
	queryValue := reflect.ValueOf(query)
	queryFields := queryType.NumField()
	fmt.Println("queryFields", queryFields, queryType, queryValue)
	for i := 0; i < queryFields; i++ {
		field := queryType.Field(i)
		value := queryValue.Field(i)
		tag := field.Tag
		column := tag.Get("column")
		operate := tag.Get("operate")
		fmt.Println("field", field, "value", value, "tag", tag, "column", column, "operate", operate)
		// 代码中对 Select、Order 和 Expand 这三个特殊字段进行了处理，但这些不是对指定字段值的过滤，
		// 而是用于指定查询的字段、排序规则和预加载关联数据。
		switch field.Name {
		case "Select":
			if selectFields := strings.Split(value.String(), ","); len(selectFields) > 1 {
				db = db.Select(selectFields)
			}
		case "Order":
			enum := tag.Get("enum")
			if value.String() != "" {
				orderFields := strings.Split(value.String(), ",")
				if len(orderFields) > 0 {
					for _, orderField := range orderFields {
						fieldArray := strings.Split(orderField, " ")
						if !strings.Contains(enum, fieldArray[0]) || len(fieldArray) != 2 {
							continue
						}
						db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: fieldArray[0]}, Desc: fieldArray[1] == "desc"})
					}
				} else {
					db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
				}
			} else {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
			}
		case "Expand":
			enum := tag.Get("enum")
			if value.String() != "" {
				expandFields := strings.Split(value.String(), ",")
				for _, expandField := range expandFields {
					if !strings.Contains(enum, expandField) {
						continue
					}
					db = db.Preload(expandField)
				}
			}
		}

		if !value.IsValid() || value.IsZero() {
			continue
		}
		if value.Kind() == reflect.Ptr && value.IsNil() {
			continue
		}
		if tag.Get("column") == "" {
			continue
		}
		if value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
		switch operate {
		case "$contains":
			db = db.Where(clause.And(clause.Like{Column: column, Value: "%" + value.String() + "%"}))
		case "$in":
			var values []interface{}
			for _, v := range strings.Split(value.String(), ",") {
				values = append(values, v)
			}
			db = db.Where(clause.IN{Column: column, Values: values})
		case "$gte":
			db = db.Where(clause.Gte{Column: column, Value: value.Interface()})
		case "$lte":
			db = db.Where(clause.Lte{Column: column, Value: value.Interface()})
		case "$gt":
			db = db.Where(clause.Gt{Column: column, Value: value.Interface()})
		case "$lt":
			db = db.Where(clause.Lt{Column: column, Value: value.Interface()})
		case "$not":
			db = db.Where(clause.Neq{Column: column, Value: value.Interface()})
		case "$isnull":
			if value.String() == "true" {
				db = db.Where(column + " IS NULL")
			} else {
				db = db.Where(column + " IS NOT NULL")
			}
		default:
			db = db.Where(clause.Eq{Column: column, Value: value.Interface()})
		}
	}
	if !queryValue.FieldByName("Order").IsValid() {
		db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: true})
	}
	return db
}

func FilterByQuery(query interface{}) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		model := db.Statement.Model
		modelType := reflect.TypeOf(model)
		if model != nil && modelType.Kind() == reflect.Ptr && modelType.Elem().Kind() == reflect.Struct {
			db = expressionByField(db, query)
		}
		return db
	}
}
