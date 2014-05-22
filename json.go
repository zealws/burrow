package burrow

import (
	"encoding/json"
	"github.com/pilu/traffic"
	"net/http"
)

func marshalObject(lm linkManager, crud CRUD, obj interface{}) ([]byte, error) {
	links, err := lm.AllLinksFor(crud, obj)
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
	x["links"] = links
	return json.Marshal(x)
}

func writeJsonObject(lm linkManager, crud CRUD, w traffic.ResponseWriter, obj interface{}) {
	bytes, err := marshalObject(lm, crud, obj)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteText("Could not marshal json response.")
	}
	w.Write(bytes)
}

func writeJsonObjects(lm linkManager, crud CRUD, w traffic.ResponseWriter, objs []interface{}) {
	// Warning: Magic Number here.
	// Want this to be as close to the size in byte of the resulting JSON object.
	// Guestimating 100 bytes per object in the list. Probably short, but it's a place to start.
	bytes := make([]byte, 0, len(objs)*100)
	bytes = append(bytes, byte('['))
	first := true
	for _, o := range objs {
		if first {
			first = false
		} else {
			bytes = append(bytes, byte(','))
		}
		byts, err := marshalObject(lm, crud, o)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.WriteText("Could not marshal json response.")
		}
		bytes = append(bytes, byts...)
	}
	bytes = append(bytes, byte(']'))
	w.Write(bytes)
}
