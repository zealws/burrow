package burrow

import (
	"encoding/json"
	"fmt"
	"github.com/pilu/traffic"
	"net/http"
	"strings"
)

/*
Format the links using the given host string, which should be `hostname:port`.
This method modifies the given map.
*/
func fixLinks(host string, links *map[string]string) {
	for name, url := range *links {
		if !strings.Contains(url, "http:") {
			(*links)[name] = fmt.Sprintf("http://%s%s", host, url)
		}
	}
}

func marshalObject(api *Api, host string, t ApiDescription, obj interface{}) ([]byte, error) {
	links, err := linksFor(api, t, obj)
	if err != nil {
		return nil, err
	}
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

func writeJsonObject(api *Api, host string, w traffic.ResponseWriter, t ApiDescription, obj interface{}) {
	bytes, err := marshalObject(api, host, t, obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteText("Could not marshal json response.")
	}
	w.Write(bytes)
}

func writeJsonObjects(api *Api, host string, w traffic.ResponseWriter, t ApiDescription, objs []interface{}) {
	bytes := make([]byte, 0, len(objs)*100)
	bytes = append(bytes, byte('['))
	first := true
	for _, o := range objs {
		if first {
			first = false
		} else {
			bytes = append(bytes, byte(','))
		}
		byts, err := marshalObject(api, host, t, o)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.WriteText("Could not marshal json response.")
		}
		bytes = append(bytes, byts...)
	}
	bytes = append(bytes, byte(']'))
	w.Write(bytes)
}
