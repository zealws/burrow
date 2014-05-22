package main

import (
	"github.com/zfjagann/burrow"
)

type Book struct {
	Id        int `rest:"id"`
	Name      string
	ISBN      string
	Author    string
	LibraryId int `references:"Library" link_name:"owner"`
}

func GetBook(id int) (interface{}, error) {
	if id < 0 || id >= len(AllBooks) {
		return nil, burrow.ApiError(404, "Could not find book with id:", id)
	}
	return &AllBooks[id], nil
}

func GetBooks() ([]interface{}, error) {
	stuff := make([]interface{}, len(AllBooks))
	for i, thing := range AllBooks {
		stuff[i] = thing
	}
	return stuff, nil
}

func UpdateBook(obj interface{}) error {
	book := obj.(*Book)
	AllBooks[book.Id] = *book
	return nil
}

type Library struct {
	Id       int `rest:"id"`
	Name     string
	Location string
}

func GetLibrary(id int) (interface{}, error) {
	if id < 0 || id >= len(AllLibraries) {
		return nil, burrow.ApiError(404, "Could not find library with id:", id)
	}
	return &AllLibraries[id], nil
}

func GetLibraries() ([]interface{}, error) {
	stuff := make([]interface{}, len(AllLibraries))
	for i, thing := range AllLibraries {
		stuff[i] = thing
	}
	return stuff, nil
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

	api.Add(burrow.New(Book{}, nil, GetBook, GetBooks, UpdateBook, nil))
	api.Add(burrow.New(Library{}, nil, GetLibrary, GetLibraries, nil, nil))
	api.Serve("0.0.0.0", 8080)
}
