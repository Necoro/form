package form

import (
	"html/template"
	"reflect"
	"strings"
)

// valueOf is basically just reflect.ValueOf, but if the Kind() of the
// value is a pointer or interface it will try to get the reflect.Value
// of the underlying element, and if the pointer is nil it will
// create a new instance of the type and return the reflect.Value of it.
//
// This is used to make the rest of the fields function simpler.
func valueOf(v interface{}) reflect.Value {
	rv := reflect.ValueOf(v)
	// If a nil pointer is passed in but has a type we can recover, but I
	// really should just panic and tell people to fix their shitty code.
	if rv.Type().Kind() == reflect.Pointer && rv.IsNil() {
		rv = reflect.Zero(rv.Type().Elem())
	}
	// If we have a pointer or interface let's try to get the underlying
	// element
	for rv.Kind() == reflect.Pointer || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv
}

func fields(v interface{}, names ...string) []field {
	rv := valueOf(v)
	if rv.Kind() != reflect.Struct {
		// We can't really do much with a non-struct type. I suppose this
		// could eventually support maps as well, but for now it does not.
		panic("invalid value; only structs are supported")
	}

	t := rv.Type()
	vFields := reflect.VisibleFields(t)
	ret := make([]field, 0, len(vFields))
	for _, tf := range vFields {
		if !tf.IsExported() {
			continue
		}

		rf := rv.FieldByIndex(tf.Index)
		// If this is a nil pointer, create a new instance of the element.
		if tf.Type.Kind() == reflect.Pointer && rf.IsNil() {
			rf = reflect.Zero(tf.Type.Elem())
		}

		// If this is a struct it has nested fields we need to add. The
		// simplest way to do this is to recursively call `fields` but
		// to provide the name of this struct field to be added as a prefix
		// to the fields.
		// This does not apply to anonymous structs, because their fields are
		// seen as "inlined".
		if reflect.Indirect(rf).Kind() == reflect.Struct {
			if !tf.Anonymous {
				ret = append(ret, fields(rf.Interface(), append(names, tf.Name)...)...)
			}
			continue
		}

		// If we are still in this loop then we aren't dealing with a nested
		// struct and need to add the field. First we check to see if the
		// ignore tag is present, then we set default values, then finally
		// we overwrite defaults with any provided tags.
		tags, ignored := parseTags(tf.Tag.Get("form"))
		if ignored {
			continue
		}
		name := append(names, tf.Name)
		f := field{
			Name:        strings.Join(name, "."),
			Label:       tf.Name,
			Placeholder: tf.Name,
			Type:        "text",
			Value:       rf.Interface(),
		}
		f.applyTags(tags)
		ret = append(ret, f)
	}
	return ret
}

func (f *field) applyTags(tags map[string]string) {
	if v, ok := tags["name"]; ok {
		f.Name = v
	}
	if v, ok := tags["label"]; ok {
		f.Label = v
		// DO NOT move this label check after the placeholder check or
		// this will cause issues.
		f.Placeholder = v
	}
	if v, ok := tags["placeholder"]; ok {
		f.Placeholder = v
	}
	if v, ok := tags["type"]; ok {
		f.Type = v
	}
	if v, ok := tags["id"]; ok {
		f.ID = v
	}
	if v, ok := tags["footer"]; ok {
		// Probably shouldn't be HTML but whatever.
		f.Footer = template.HTML(v)
	}
	if v, ok := tags["class"]; ok {
		f.Class = v
	}
	if v, ok := tags["readonly"]; ok {
		f.ReadOnly = v == "true"
	}
}

func parseTags(tags string) (map[string]string, bool) {
	tags = strings.TrimSpace(tags)
	if len(tags) == 0 {
		return map[string]string{}, false
	}
	split := strings.Split(tags, ";")
	ret := make(map[string]string, len(split))
	for _, tag := range split {
		kv := strings.Split(tag, "=")
		if len(kv) < 2 {
			if kv[0] == "-" {
				return nil, true
			}
			continue
		}
		k, v := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])
		ret[k] = v
	}
	return ret, false
}

type field struct {
	Name        string
	Label       string
	Placeholder string
	Type        string
	ID          string
	ReadOnly    bool
	Value       interface{}
	Footer      template.HTML
	Class       string
}
