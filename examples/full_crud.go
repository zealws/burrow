package main

import (
	"github.com/zfjagann/burrow"
)

type Book struct {
	Id     int `rest:"id"`
	Name   string
	ISBN   string
	Author string
}

var NextId int = 0
var AllBooks map[int]Book

func CreateBook() (interface{}, error) {
	defer func() { NextId++ }()
	AllBooks[NextId] = Book{NextId, "", "", ""}
	return AllBooks[NextId], nil
}

func GetBook(id int) (interface{}, error) {
	if id < 0 || id >= len(AllBooks) {
		return nil, burrow.ApiError(404, "Could not find book with id:", id)
	}
	return AllBooks[id], nil
}

func GetBooks() ([]interface{}, error) {
	stuff := make([]interface{}, len(AllBooks))
	for i, thing := range AllBooks {
		stuff[i] = thing
	}
	return stuff, nil
}

func UpdateBook(obj interface{}) error {
	book := obj.(Book)
	AllBooks[book.Id] = book
	return nil
}

func DeleteBook(id int) error {
	_, ok := AllBooks[id]
	if !ok {
		return burrow.ApiError(404, "Could not find book with id", id)
	}
	delete(AllBooks, id)
	return nil
}

func init() {
	AllBooks = make(map[int]Book)
	AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens"}
	AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe"}
	AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare"}
	NextId = 3
}

func main() {
	api := burrow.NewApi()

	// Create a CRUD which is used to manage our Book object
	bookCrud := burrow.New(Book{}, CreateBook, GetBook, GetBooks, UpdateBook, nil)

	// Add the Book CRUD to the server
	api.Add(bookCrud)

	// Run the server!
	api.Serve("0.0.0.0", 8080)
}
