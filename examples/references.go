package main

import (
	"github.com/zfjagann/burrow"
)

type Book struct {
	Id        int `rest:"id"`
	Name      string
	ISBN      string
	Author    string
	LibraryId int `references:"library" link_name:"owner"`
}

func GetBook(id int) interface{} {
	return &AllBooks[id]
}

func GetBooks() []interface{} {
	stuff := make([]interface{}, len(AllBooks))
	for i, thing := range AllBooks {
		stuff[i] = thing
	}
	return stuff
}

type Library struct {
	Id       int `rest:"id"`
	Name     string
	Location string
}

func GetLibrary(id int) interface{} {
	return &AllLibraries[id]
}

func GetLibraries() []interface{} {
	stuff := make([]interface{}, len(AllLibraries))
	for i, thing := range AllLibraries {
		stuff[i] = thing
	}
	return stuff
}

var AllBooks []Book
var AllLibraries []Library

func init() {
	AllBooks = make([]Book, 3)
	AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens", 0}
	AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe", 0}
	AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare", 1}

	AllLibraries = make([]Library, 2)
	AllLibraries[0] = Library{0, "Mountain View Public Library", "Mountain View"}
	AllLibraries[1] = Library{1, "Cupertino Public Library", "Cupertino"}
}

func main() {
	api := burrow.NewApi()
	api.AddApi(Book{}, GetBook, GetBooks)
	api.AddApi(Library{}, GetLibrary, GetLibraries)
	api.Serve("0.0.0.0", 8080)
}
