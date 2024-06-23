package stdlib

import "reflect"

// AnyTo converts an any interface value to a value
// that can be type asserted to type T.
func AnyTo[T any](v any) any {
	rt := reflect.TypeFor[T]()
	vt := reflect.TypeOf(v)

	if rt.Kind() == reflect.Ptr {
		// t = *thing; v = thing
		if vt.Kind() == reflect.Ptr {
			return v
		} else {
			// t = *thing; v = thing
			return &v
		}
	} else {
		// t = thing; v = *thing
		if vt.Kind() == reflect.Ptr {
			return *v.(*T)
		} else {
			// t = thing; v = thing
			return v
		}
	}
}
