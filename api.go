package rest

import (
	"encoding/json"
	"fmt"
	"github.com/pilu/traffic"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type SpecificGetter func(int) interface{}
type GeneralGetter func() []interface{}

type ApiType interface {
	Name() string
	Reflect() reflect.Type
	RootGetter() GeneralGetter
	IndividualGetter() SpecificGetter
}

type typeImpl struct {
	name string
	refl reflect.Type
	spec SpecificGetter
	gen  GeneralGetter
}

func (t typeImpl) Name() string {
	return t.name
}

func (t typeImpl) Reflect() reflect.Type {
	return t.refl
}

func (t typeImpl) IndividualGetter() SpecificGetter {
	return t.spec
}

func (t typeImpl) RootGetter() GeneralGetter {
	return t.gen
}

/*
ApiType Helpers
*/

func getApiType(obj interface{}, spec SpecificGetter, gen GeneralGetter) ApiType {
	typ := reflect.TypeOf(obj)
	return typeImpl{
		name: typ.Name(),
		refl: typ,
		spec: spec,
		gen:  gen,
	}
}

func getIdField(t ApiType) string {
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

func getIdValue(t ApiType, obj interface{}) (int, bool) {
	idField := getIdField(t)
	if idField == "" {
		return 0, false
	}
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = reflect.Indirect(val)
	}
	fld := val.FieldByName(idField)
	return int(fld.Int()), true
}

func selfLinkFor(t ApiType, obj interface{}) (self string, ok bool) {
	id, ok := getIdValue(t, obj)
	self = strings.ToLower(fmt.Sprintf("/%s/%d", t.Name(), id))
	return
}

func indexLinkFor(t ApiType) string {
	return strings.ToLower("/" + t.Name())
}

func linksFor(t ApiType, obj interface{}) map[string]string {
	links := make(map[string]string)
	links["index"] = indexLinkFor(t)
	self, ok := selfLinkFor(t, obj)
	if ok {
		links["self"] = self
	}
	return links
}

/*
Handler Generators
*/

func makeSpecificHandler(t ApiType) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintf(w, "Invalid ID: ", params.Get("id"))
		}
		writeJsonObject(r.Host, w, t, t.IndividualGetter()(id))
	}
	return traffic.HttpHandleFunc(f)
}

func makeGenericHandler(t ApiType) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		writeJsonObjects(r.Host, w, t, t.RootGetter()())
	}
	return traffic.HttpHandleFunc(f)
}

func getHandlers(t ApiType) map[string]traffic.HttpHandleFunc {
	results := make(map[string]traffic.HttpHandleFunc)
	root := indexLinkFor(t)
	generic := makeGenericHandler(t)
	results[root] = generic
	results[root+"/"] = generic
	results[root+"/:id"] = makeSpecificHandler(t)
	return results
}

/*
JSON FormatterswriteJsonObjects
*/

func fixLinks(host string, links *map[string]string) {
	for name, url := range *links {
		(*links)[name] = fmt.Sprintf("http://%s%s", host, url)
	}
}

func marshalObject(host string, t ApiType, obj interface{}) ([]byte, error) {
	links := linksFor(t, obj)
	jsn, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	x := make(map[string]interface{})
	err = json.Unmarshal(jsn, &x)
	if err != nil {
		return nil, err
	}
	fixLinks(host, &links)
	x["links"] = links
	return json.Marshal(x)
}

func writeJsonObject(host string, w traffic.ResponseWriter, t ApiType, obj interface{}) {
	bytes, err := marshalObject(host, t, obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteText("Could not marshal json response.")
	}
	w.Write(bytes)
}

func writeJsonObjects(host string, w traffic.ResponseWriter, t ApiType, objs []interface{}) {
	bytes := make([]byte, 0, len(objs)*100)
	bytes = append(bytes, byte('['))
	first := true
	for _, o := range objs {
		if first {
			first = false
		} else {
			bytes = append(bytes, byte(','))
		}
		byts, err := marshalObject(host, t, o)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.WriteText("Could not marshal json response.")
		}
		bytes = append(bytes, byts...)
	}
	bytes = append(bytes, byte(']'))
	w.Write(bytes)
}

/*
Api Type
*/

type Api struct {
	types []ApiType
}

func NewApi() (api *Api) {
	return &Api{make([]ApiType, 0)}
}

func (api *Api) AddApi(obj interface{}, spec SpecificGetter, gen GeneralGetter) {
	typ := getApiType(obj, spec, gen)
	api.types = append(api.types, typ)
}

func (api *Api) Serve(addr string, port int) {
	traffic.SetPort(port)
	traffic.SetHost(addr)
	router := traffic.New()
	for _, t := range api.types {
		for url, handler := range getHandlers(t) {
			router.Get(url, handler)
		}
	}
	router.Run()
}
