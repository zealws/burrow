package burrow

import (
	"fmt"
	"github.com/pilu/traffic"
	"net/http"
	"strconv"
	"strings"
)

/*
Api Type
*/

type Api struct {
	types []ApiDescription
}

func NewApi() (api *Api) {
	return &Api{make([]ApiDescription, 0)}
}

func (api *Api) AddApi(obj interface{}, getter ObjectGetter, index IndexGetter) {
	typ := getApiDescription(obj, getter, index)
	api.types = append(api.types, typ)
}

func (api *Api) Add(typ ApiDescription) {
	api.types = append(api.types, typ)
}

func (api *Api) Serve(addr string, port int) {
	traffic.SetPort(port)
	traffic.SetHost(addr)
	router := traffic.New()
	for _, t := range api.types {
		api.addHandlersTo(t, router)
	}
	router.Get("/", api.rootHandler())
	router.Run()
}

func (api *Api) makeIndexLinks() map[string]string {
	links := make(map[string]string)
	for _, t := range api.types {
		links[t.Name()+" index"] = indexLinkFor(t)
	}
	links["root"] = "/"
	links["self"] = "/"
	return links
}

func (api *Api) rootHandler() traffic.HttpHandleFunc {
	links := api.makeIndexLinks()
	handler := func(w traffic.ResponseWriter, r *traffic.Request) {
		root := make(map[string]interface{})
		// Have to copy the map each time since fixLinks modifies it, and the links object
		// above is common to all requests.
		// Not a huge deal since it's O(Number of API object types), which should be small.
		myLinks := copyMap(links)
		fixLinks(r.Host, &myLinks)
		root["links"] = myLinks
		w.WriteJSON(root)
	}
	return traffic.HttpHandleFunc(handler)
}

func copyMap(old map[string]string) map[string]string {
	knew := make(map[string]string)
	for k, v := range old {
		knew[k] = v
	}
	return knew
}

func (api *Api) makeSpecificHandler(t ApiDescription) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		params := r.URL.Query()
		id, err := strconv.Atoi(params.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Fprintln(w, "Invalid ID: ", params.Get("id"))
		}
		writeJsonObject(api, r.Host, w, t, t.Getter()(id))
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) makeGenericHandler(t ApiDescription) traffic.HttpHandleFunc {
	f := func(w traffic.ResponseWriter, r *traffic.Request) {
		writeJsonObjects(api, r.Host, w, t, t.Index()())
	}
	return traffic.HttpHandleFunc(f)
}

func (api *Api) addHandlersTo(t ApiDescription, router *traffic.Router) map[string]traffic.HttpHandleFunc {
	results := make(map[string]traffic.HttpHandleFunc)
	root := indexLinkFor(t)
	generic := api.makeGenericHandler(t)
	router.Get(root, generic)
	router.Get(root+"/", generic)
	router.Get(root+"/:id", api.makeSpecificHandler(t))
	return results
}

func (api *Api) getType(name string) *ApiDescription {
	for _, t := range api.types {
		if strings.ToLower(name) == strings.ToLower(t.Name()) {
			return &t
		}
	}
	return nil
}
