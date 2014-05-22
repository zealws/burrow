package burrow

import (
	"errors"
	"fmt"
	"strings"
)

type linkManager struct {
	cruds    []CRUD
	hostspec string
}

func (m linkManager) SpecificLink(typeName string, id int) (string, error) {
	crud, err := m.findCrud(typeName)
	if err != nil {
		return "", err
	}
	return m.ObjectLink(crud, id), nil
}

func (m linkManager) SelfLink(crud CRUD, obj interface{}) (string, bool) {
	id, ok := getIdValue(crud, obj)
	self := m.ObjectLink(crud, id)
	return self, ok
}

func (m linkManager) ObjectLink(crud CRUD, id int) string {
	return m.link(fmt.Sprintf("/%s/%d", crud.Name(), id))
}

func (m linkManager) IndexLink(crud CRUD) string {
	return m.link("/" + crud.Name())
}

func (m linkManager) IndexUrl(crud CRUD) string {
	return m.url("/" + crud.Name())
}

func (m linkManager) RootLinks() map[string]string {
	links := make(map[string]string)
	for _, c := range m.cruds {
		links[c.Name()+" index"] = m.IndexLink(c)
	}
	links["root"] = "/"
	links["self"] = "/"
	return links
}

func (m linkManager) findCrud(name string) (CRUD, error) {
	for _, c := range m.cruds {
		if c.Name() == name {
			return c, nil
		}
	}
	return nil, errors.New("No CRUD exists named " + name)
}

func (m linkManager) link(partial string) string {
	return fmt.Sprintf("http://%s%s", m.hostspec, m.url(partial))
}

func (m linkManager) url(partial string) string {
	return strings.ToLower(partial)
}

func (m linkManager) AllLinksFor(c CRUD, obj interface{}) (map[string]string, error) {
	links := make(map[string]string)
	links[c.Name()+" index"] = m.IndexLink(c)
	links["root"] = "/"
	self, ok := m.SelfLink(c, obj)
	if ok {
		links["self"] = self
	}
	fields, err := linkedFieldsFor(c, obj)
	if err != nil {
		return nil, err
	}
	for _, ref := range fields {
		link, err := m.SpecificLink(ref.typeName, ref.id)
		if err != nil {
			return nil, err
		}
		links[ref.name] = link
	}
	return links, err
}
