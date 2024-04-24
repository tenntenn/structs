package structs

import (
	"cmp"
	"fmt"
	"go/ast"
	"reflect"
	"slices"
)

type Field struct {
	typ reflect.StructField
	val reflect.Value
}

type FieldOption func(*Field)

func Tag[Tag string | reflect.StructTag](tag Tag) FieldOption {
	return func(f *Field) {
		f.typ.Tag = reflect.StructTag(tag)
	}
}

func F(nameAndValue ...any) *Field {

	if len(nameAndValue) < 2 {
		panic("name and value must be specified")
	}

	name, ok := nameAndValue[0].(string)
	if !ok {
		panic(fmt.Sprintf("invalid name: %[1]v(%[1]T)", nameAndValue[0]))
	}

	if !ast.IsExported(name) {
		panic(fmt.Sprintf("the field must be exported: %s", name))
	}

	v := reflect.ValueOf(nameAndValue[1])

	typ := reflect.StructField{
		Name: name,
		Type: v.Type(),
	}

	field := &Field{
		typ: typ,
		val: v,
	}

	for _, opt := range nameAndValue[2:] {
		opt, ok := opt.(FieldOption)
		if ok {
			opt(field)
		}
	}

	return field
}

// Merge make new struct value with the exported fields of the given structs.
// If some have a same name of field, the field of the last struct in the parameters is used.
// If non struct value is in the parameters, it would be ignored.
//
//	type A { N int `json:"n"`}
//	type B { S string `json:"s"`}
//
//	/*
//	struct{
//	     N int	`json:"n"`
//	     S string	`json:"s"`
//	}{
//		N: 100,
//		S: "sample",
//	}
//	*/
//	Merge(A{N:100}, B{S:"sample"})
func Merge(xs ...any) any {
	typeAndValues := make(map[string]*Field)
	for _, x := range xs {
		putFieldsTo(typeAndValues, x)
	}

	return newWith(typeAndValues)
}

// Of make new struct value with the given names, values and tags as fields.
//
//	/*
//	struct{
//	     N int
//	     S string
//	}{
//		N: 100,
//		S: "sample",
//	}
//	*/
//	Of(F("N", 100), F("S", "sample"))
//
//	/*
//	struct{
//	     N int	`json:"n"`
//	     S string	`json:"s"`
//	}{
//		N: 100,
//		S: "sample",
//	}
//	*/
//	Of(F("N", 100, Tag(`json:"n"`)), F("S", "sample", Tag(`json:"s"`)))
//
// The tags can be ommited.
func Of(fields ...*Field) any {

	if len(fields) == 0 {
		return struct{}{}
	}

	typeAndValues := make(map[string]*Field, len(fields))
	for _, field := range fields {
		typeAndValues[field.typ.Name] = field
	}

	return newWith(typeAndValues)
}

func putFieldsTo(fields map[string]*Field, x any) {
	xv := reflect.ValueOf(x)
	if xv.Kind() == reflect.Ptr {
		xv = xv.Elem()
	}

	if xv.Kind() != reflect.Struct {
		return
	}

	xt := xv.Type()

	for i := range xt.NumField() {
		typ := xt.Field(i)
		if !typ.IsExported() {
			continue
		}

		fields[typ.Name] = &Field{
			typ: typ,
			val: xv.Field(i),
		}
	}
}

func newWith(typeAndValues map[string]*Field) any {
	fields := make([]reflect.StructField, 0, len(typeAndValues))
	for _, tv := range typeAndValues {
		fields = append(fields, tv.typ)
	}

	slices.SortFunc(fields, func(a, b reflect.StructField) int {
		return cmp.Compare(a.Name, b.Name)
	})

	ptr := reflect.New(reflect.StructOf(fields))
	v := ptr.Elem()

	for name, tv := range typeAndValues {
		f := v.FieldByName(name)
		if f.CanSet() && tv.val.Type().AssignableTo(f.Type()) {
			f.Set(tv.val)
		}
	}

	return v.Interface()
}
