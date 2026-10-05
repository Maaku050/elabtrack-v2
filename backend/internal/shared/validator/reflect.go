package validator

import "reflect"

// reflectJSONTag reads the `json` struct tag for a field, returning the
// name portion (before any ",omitempty").
func reflectJSONTag(s any, field string) string {
	t := reflect.TypeOf(s)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return ""
	}
	f, ok := t.FieldByName(field)
	if !ok {
		return ""
	}
	tag := f.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	// strip ",omitempty" etc.
	for i := 0; i < len(tag); i++ {
		if tag[i] == ',' {
			return tag[:i]
		}
	}
	return tag
}
