package main

import (
	"github.com/zfjagann/burrow"
)

/*
The API Object we want to expose.
*/
type Book struct {
	/* This tag allow the "self" link to be generated. */
	Id     int `rest:"id"`
	Name   string
	ISBN   string
	Author string
}

var AllBooks []Book

/*
Accessor for a single book, by id.

This function would likely be a database lookup in a more realistic application.
*/
func GetBook(id int) (interface{}, error) {
	if id < 0 || id >= len(AllBooks) {
		return nil, burrow.ApiError(404, "Could not find book with id:", id)
	}
	return &AllBooks[id], nil
}

/*
Index of all the books.

This function would likely be a database lookup in a more realistic application.
*/
func GetBooks() ([]interface{}, error) {
	stuff := make([]interface{}, len(AllBooks))
	for i, thing := range AllBooks {
		stuff[i] = thing
	}
	return stuff, nil
}

/*
Update a given book.
*/
func UpdateBook(obj interface{}) error {
	book := obj.(*Book)
	AllBooks[book.Id] = *book
	return nil
}

/*
Load our test data into memory.
*/
func init() {
	AllBooks = make([]Book, 3)
	AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens"}
	AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe"}
	AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare"}
}

func main() {
	api := burrow.NewApi()

	// Create a CRUD which is used to manage our Book object
	bookCrud := burrow.New(Book{}, nil, GetBook, GetBooks, UpdateBook, nil)

	// Add the Book CRUD to the server
	api.Add(bookCrud)

	// Run the server!
	api.Serve("0.0.0.0", 8080)
}
