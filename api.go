package burrow

import (
	"fmt"
	"github.com/pilu/traffic"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

/**
*** Error Types
***/

type Error struct {
	message string
	status  int
}

func (e Error) Error() string {
	return e.message
}

func NewError(items ...interface{}) Error {
	return Error{fmt.Sprint(items...), 500}
}

func ApiError(status int, items ...interface{}) Error {
	return Error{fmt.Sprint(items...), status}
}

func handle(w traffic.ResponseWriter, e error) {
	switch err := e.(type) {
	case Error:
		w.WriteHeader(err.status)
		fmt.Fprintln(w, err)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintln(w, err)
	}

}

/**
*** Api Type
**/

type Api struct {
	types []CRUD
}

func NewApi() (api *Api) {
	return &Api{make([]CRUD, 0)}
}

func (api *Api) Add(crud CRUD) {
	api.types = append(api.types, crud)
}

func (api *Api) Serve(addr string, port int) {
	lm := linkManager{api.types, ""}
	traffic.SetPort(port)
	traffic.SetHost(addr)
	router := traffic.New()
	for _, t := range api.types {
		api.addHandlersTo(lm, t, router)
	}
	router.Get("/", api.rootHandler())
	router.Run()
}

func (api *Api) rootHandler() traffic.HttpHandleFunc {
	handler := func(w traffic.ResponseWriter, r *traffic.Request) {
		lm := api.linkManager(r)
		root := make(map[string]interface{})
		root["links"] = lm.RootLinks()
		w.WriteJSON(root)
	}
	return traffic.HttpHandleFunc(handler)
}

func (api *Api) makeReadHandler(crud CRUD) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		lm := api.linkManager(r)
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			handle(w, err)
			return
		}
		obj, err := crud.Read(id)
		if err != nil {
			handle(w, err)
			return
		}
		err = writeJsonObject(lm, crud, w, obj)
		if err != nil {
			handle(w, err)
			return
		}
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) makeIndexHandler(crud CRUD) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		lm := api.linkManager(r)
		index, err := crud.Index()
		if err != nil {
			handle(w, err)
			return
		}
		err = writeJsonObjects(lm, crud, w, index)
		if err != nil {
			handle(w, err)
			return
		}
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) makeCreateHandler(crud CRUD) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		lm := api.linkManager(r)
		obj, err := crud.Create()
		if err != nil {
			handle(w, err)
			return
		}
		w.WriteHeader(http.StatusCreated)
		err = writeJsonObject(lm, crud, w, obj)
		if err != nil {
			handle(w, err)
			return
		}
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) makeUpdateHandler(crud CRUD) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		lm := api.linkManager(r)
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			handle(w, err)
			return
		}
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handle(w, err)
			return
		}
		obj, err := crud.Read(id)
		if err != nil {
			handle(w, err)
			return
		}
		obj, err = updateObject(crud, obj, bytes)
		if err != nil {
			handle(w, err)
			return
		}
		err = crud.Update(obj)
		if err != nil {
			handle(w, err)
			return
		}
		err = writeJsonObject(lm, crud, w, obj)
		if err != nil {
			handle(w, err)
			return
		}
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) makeDeleteHandler(crud CRUD) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			handle(w, err)
			return
		}
		err = crud.Delete(id)
		if err != nil {
			handle(w, err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) addHandlersTo(lm linkManager, crud CRUD, router *traffic.Router) map[string]traffic.HttpHandleFunc {
	results := make(map[string]traffic.HttpHandleFunc)
	root := lm.IndexUrl(crud)
	router.Get(root, api.makeIndexHandler(crud))
	router.Post(root, api.makeCreateHandler(crud))
	router.Put(root+"/:id", api.makeUpdateHandler(crud))
	router.Delete(root+"/:id", api.makeDeleteHandler(crud))
	router.Get(root+"/:id", api.makeReadHandler(crud))
	return results
}

func (api *Api) getCRUD(name string) CRUD {
	for _, t := range api.types {
		if strings.ToLower(name) == strings.ToLower(t.Name()) {
			return t
		}
	}
	return nil
}

func (api *Api) linkManager(r *traffic.Request) linkManager {
	return linkManager{api.types, r.Host}
}
