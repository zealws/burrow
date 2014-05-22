package burrow

import (
	"encoding/json"
	"errors"
	"reflect"
)

/**
*** CRUD Interface
**/

/*
Represents the interface of interaction between the backend Object store and the API.
*/
type CRUD interface {
	/*
		Returns the name of the type that this CRUD represents.

		This is used to determine URIs.
	*/
	Name() string

	/*
		Returns the reflected type of the object this CRUD manages.
	*/
	Reflect() reflect.Type

	/*
		Create a new object and return it.
	*/
	Create() (interface{}, error)

	/*
		Update a given object.
		Returns nil if successful, and an error otherwise.
	*/
	Update(interface{}) error

	/*
		Load and return a single object.
		The error return value will be nil if successful, and non-nil otherwise.
	*/
	Read(int) (interface{}, error)

	/*
		Load and return all objects.
		The error return value will be nil if succesful, and non-nil otherwise.
	*/
	Index() ([]interface{}, error)

	/*
		Delete the object specified by the given id.
		Returns nil if successful, and an error otherwise.
	*/
	Delete(int) error
}

/**
*** CRUD Helper Methods
**/

type reference struct {
	typeName string
	name     string
	id       int
}

func getIdField(c CRUD) string {
	typ := c.Reflect()
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.FieldByIndex([]int{i})
		restTag := fld.Tag.Get("rest")
		if restTag == "id" {
			return fld.Name
		}
	}
	return ""
}

func getIdValue(c CRUD, obj interface{}) (int, bool) {
	idField := getIdField(c)
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

func updateObject(c CRUD, obj interface{}, body []byte) (err error) {
	fields := make(map[string]interface{})
	err = json.Unmarshal(body, &fields)
	if err != nil {
		return
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typ := val.Type()
	for name, value := range fields {
		_, ok := typ.FieldByName(name)
		if !ok {
			return ApiError(406, "Could not find", c.Name(), "field", name)
		}
		val.FieldByName(name).Set(reflect.ValueOf(value))
	}
	return nil
}

func linkedFieldsFor(c CRUD, obj interface{}) ([]reference, error) {
	linkedFields := make([]reference, 0)
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		fld := typ.FieldByIndex([]int{i})
		tag := fld.Tag
		ref := tag.Get("references")
		if ref != "" {
			name := tag.Get("link_name")
			if name == "" {
				name = ref
			}
			// Check that runtime type is a permissible ID type.
			if !validIdValue(val.FieldByIndex([]int{i})) {
				return nil, errors.New("Reference of non-integer type: " + c.Name() + "." + fld.Name + ". Check to be sure your API models are defined correctly.")
			}
			id := val.FieldByIndex([]int{i}).Int()
			linkedFields = append(linkedFields, reference{ref, name, int(id)})
		}
	}
	return linkedFields, nil
}

/**
*** Concrete CRUD type below.
**/

type crudImpl struct {
	typ     reflect.Type
	creator func() (interface{}, error)
	updater func(interface{}) error
	reader  func(int) (interface{}, error)
	indexer func() ([]interface{}, error)
	deleter func(interface{}) error
}

func New(
	value interface{},
	creator func() (interface{}, error),
	reader func(int) (interface{}, error),
	indexer func() ([]interface{}, error),
	updater func(interface{}) error,
	deleter func(interface{}) error,
) CRUD {
	return &crudImpl{reflect.TypeOf(value), creator, updater, reader, indexer, deleter}
}

func (crud *crudImpl) Name() string {
	return crud.typ.Name()
}

func (crud *crudImpl) Reflect() reflect.Type {
	return crud.typ
}

func (crud *crudImpl) Create() (interface{}, error) {
	if crud.creator == nil {
		return nil, errors.New("Creation of " + crud.Name() + " is not allowed.")
	}
	return crud.creator()
}

func (crud *crudImpl) Update(obj interface{}) error {
	if crud.updater == nil {
		return errors.New("Updates of " + crud.Name() + " are not allowed.")
	}
	return crud.updater(obj)
}

func (crud *crudImpl) Read(id int) (interface{}, error) {
	if crud.reader == nil {
		return nil, errors.New("Reading " + crud.Name() + " is not allowed.")
	}
	return crud.reader(id)
}

func (crud *crudImpl) Index() ([]interface{}, error) {
	if crud.indexer == nil {
		return nil, errors.New("Reading " + crud.Name() + " is not allowed.")
	}
	return crud.indexer()
}

func (crud *crudImpl) Delete(id int) error {
	if crud.deleter == nil {
		return errors.New("Deletion of " + crud.Name() + " is not allowed.")
	}
	return crud.deleter(id)
}

/**
*** Helper Functions
**/

func validIdValue(val reflect.Value) bool {
	return val.Kind() == reflect.Int || val.Kind() == reflect.Int32 || val.Kind() == reflect.Int64
}
