package utils

import "reflect"

// ReflectType Get the name of a struct instance.
func ReflectType(obj interface{}) string {
	if t := reflect.TypeOf(obj); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

// UnpackArray Unpack an `arr` argument `any` but is exactly a type `[]any`
// Eg: We have `args` type `any`.
// var args any
// args = []string{"my-file.pdf", "png"}
// utils.UnpackArray(args)
func UnpackArray(arr any) []any {
	v := reflect.ValueOf(arr)
	r := make([]any, v.Len())

	for i := 0; i < v.Len(); i++ {
		r[i] = v.Index(i).Interface()
	}

	return r
}

// UnpackArrayT Unpack an `a` argument `any` but is exactly a type `[]T`
// Eg: We have `args` type `any`.
// var args any
// args = []string{"my-file.pdf", "png"}
// utils.UnpackArrayT[string](args)
func UnpackArrayT[T any](arr any) (r []any) {
	o := arr.([]T)
	r = make([]any, len(o))

	for i, v := range o {
		r[i] = any(v)
	}

	return r
}
