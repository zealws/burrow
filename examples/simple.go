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
func GetBook(id int) interface{} {
	return &AllBooks[id]
}

/*
Accessor for all the books.

This function would likely be a database lookup in a more realistic application.
*/
func GetBooks() []interface{} {
	stuff := make([]interface{}, len(AllBooks))
	for i, thing := range AllBooks {
		stuff[i] = thing
	}
	return stuff
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

	// Add the Book type as an API Object using the given accessors.
	api.AddApi(Book{}, GetBook, GetBooks)

	// Run the server!
	api.Serve("0.0.0.0", 8080)
}
