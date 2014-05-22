package burrow

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/*
The ObjectGetter type represents a function that is used to fetch a single object for an API.
*/
type ObjectGetter func(int) interface{}

/*
The IndexGetter type represents a function that can be used to fetch a list of objects for an API.
*/
type IndexGetter func() []interface{}

/*
The ApiDescription type includes all the information needed to fully describe a Restful API object.
*/
type ApiDescription interface {
	/*
		Returns the name of the object.
		This is used to form the URLs for this API object.
	*/
	Name() string

	/*
		Returns the reflected type of the object.
		This is used to determine run-time information about the API object.
	*/
	Reflect() reflect.Type

	/*
		Return the IndexGetter used to fetch all objects of the given type.
		This is used by API endpoint handlers to fetch all of the given API objects.
	*/
	Index() IndexGetter

	/*
		The Getter function returns the ObjectGetter type used to fetch a specific object.
		This is used by API endpoint handlers to fetch specific API objects.
	*/
	Getter() ObjectGetter
}

/*
Non-exported things below:
*/

type typeImpl struct {
	name string
	refl reflect.Type
	spec ObjectGetter
	gen  IndexGetter
}

func (t typeImpl) Name() string {
	return t.name
}

func (t typeImpl) Reflect() reflect.Type {
	return t.refl
}

func (t typeImpl) Getter() ObjectGetter {
	return t.spec
}

func (t typeImpl) Index() IndexGetter {
	return t.gen
}

func getApiDescription(obj interface{}, spec ObjectGetter, gen IndexGetter) ApiDescription {
	typ := reflect.TypeOf(obj)
	return typeImpl{
		name: strings.ToLower(typ.Name()),
		refl: typ,
		spec: spec,
		gen:  gen,
	}
}

func getIdField(t ApiDescription) string {
	typ := t.Reflect()
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.FieldByIndex([]int{i})
		restTag := fld.Tag.Get("rest")
		if restTag == "id" {
			return fld.Name
		}
	}
	return ""
}

func getIdValue(t ApiDescription, obj interface{}) (int, bool) {
	idField := getIdField(t)
	if idField == "" {
		return 0, false
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	fld := val.FieldByName(idField)
	// Try to protect from panic if the field is not actually an integer
	if !validIdValue(fld) {
		return 0, false
	}
	return int(fld.Int()), true
}

func selfLinkFor(t ApiDescription, obj interface{}) (self string, ok bool) {
	id, ok := getIdValue(t, obj)
	self = linkFor(t, id)
	return
}

func linkFor(t ApiDescription, id int) string {
	return strings.ToLower(fmt.Sprintf("/%s/%d", t.Name(), id))
}

func indexLinkFor(t ApiDescription) string {
	return strings.ToLower("/" + t.Name())
}

func linksFor(api *Api, t ApiDescription, obj interface{}) (map[string]string, error) {
	links := make(map[string]string)
	links[t.Name()+" index"] = indexLinkFor(t)
	links["root"] = "/"
	self, ok := selfLinkFor(t, obj)
	if ok {
		links["self"] = self
	}
	err := addFieldLinksTo(api, t, obj, links)
	return links, err
}

func validIdValue(val reflect.Value) bool {
	return val.Kind() == reflect.Int || val.Kind() == reflect.Int32 || val.Kind() == reflect.Int64
}

func addFieldLinksTo(api *Api, t ApiDescription, obj interface{}, links map[string]string) error {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.FieldByIndex([]int{i})
		tag := fld.Tag
		reference := tag.Get("references")
		if reference != "" {
			name := tag.Get("link_name")
			if name == "" {
				name = reference
			}
			// Check that runtime type is a permissible ID type.
			if !validIdValue(val.FieldByIndex([]int{i})) {
				return errors.New("Reference of non-integer type: " + t.Name() + "." + fld.Name + ". Check to be sure your API models are defined correctly.")
			}
			id := val.FieldByIndex([]int{i}).Int()
			links[name] = linkFor(*api.getType(reference), int(id))
		}
	}
	return nil
}
