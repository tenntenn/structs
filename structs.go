package structs

import (
	"cmp"
	"fmt"
	"go/ast"
	"reflect"
	"slices"
)

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
	typeAndValues := make(map[string]field)
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
//	Of("N", 100, "S", "sample")
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
//	Of("N", 100, `json:"n"`, "S", "sample", `json:"s"`)
//
// The tags can be ommited.
func Of(nameAndValues ...any) any {

	if len(nameAndValues) == 0 {
		return struct{}{}
	}

	withTag := len(nameAndValues)%3 == 0

	if !withTag && len(nameAndValues)%2 != 0 {
		panic("invalid name and value pair")
	}

	n := 2
	if withTag {
		n = 3
	}

	typeAndValues := make(map[string]field, len(nameAndValues)/2)
	for i := 0; i < len(nameAndValues); i += n {
		name, ok := nameAndValues[i].(string)
		if !ok {
			panic(fmt.Sprintf("invalid name: %v", nameAndValues[i]))
		}

		if !ast.IsExported(name) {
			panic(fmt.Sprintf("field must be exported: %s", name))
		}

		v := reflect.ValueOf(nameAndValues[i+1])

		typ := reflect.StructField{
			Name: name,
			Type: v.Type(),
		}

		if withTag {
			tag, ok := nameAndValues[i+2].(string)
			if !ok {
				panic(fmt.Sprintf("invalid tag: %v", nameAndValues[i+2]))
			}
			typ.Tag = reflect.StructTag(tag)
		}

		typeAndValues[name] = field{
			typ: typ,
			val: v,
		}
	}

	return newWith(typeAndValues)
}

type field struct {
	typ reflect.StructField
	val reflect.Value
}

func putFieldsTo(fields map[string]field, x any) {
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

		fields[typ.Name] = field{
			typ: typ,
			val: xv.Field(i),
		}
	}
}

func newWith(typeAndValues map[string]field) any {
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
