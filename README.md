# Burrow

Burrow is a REST API Framework designed with the following in mind:

- Convention over Configuration
- Simplicity
- Ease of Use
- Written in and for [Golang](http://golang.org)

## Intent

Most REST API frameworks allow API designers to specify very complex URI schemes to represent their API Objects.
This approach to API frameworks has several drawbacks:

- The URI schemes require copious documentation since they are difficult to understand intuitively.
- URI schemes for each object must be created independently, introducing complexity and differences between similar APIs. 
- URI schemes are difficult to represent internal to the program, resulting in little to no linking between objects in an API hierarchy.

Burrow attempts to solve these by offloading the effort of defining URI schemes to the framework. This allows the
developers to write API Objects without worrying about how the resulting API will look.

This allows the framework to do the following **automatically**:

- Generate API endpoints
- Generate documentation for your API endpoints *(coming soon)*
- Manage links between referenced objects
- Generate links in your API objects to make navigating your API dead simple

## Introduction

Burrow needs three things to function properly:

- A type
- An index accessor for the type
- A specific accessor for the type

The type itself is fairly straightforward:

The Index Accessor's purpose is to fetch a list of records when asked. It should be of type `func() []interface{}`
and should return a list of objects of the given type.

The Specific Accessor's purpose is to fetch a single record based on a numerical id (types `int`,`int64`,`int32`).
It should be of type: `func(int) interface{}` and return the object referenced by the given int, and `nil`
if no such object exists.


Simple example:

```go
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

var AllBooks []Book

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

func init() {
    AllBooks = make([]Book, 3)
    AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens"}
    AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe"}
    AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare"}
}

func main() {
    api := burrow.NewApi()

    api.AddApi(Book{}, GetBook, GetBooks)

    api.Serve("0.0.0.0", 8080)
}
```

The full code (complete with comments) can be found
[here](http://github.com/zfjagann/burrow/tree/master/examples/simple.go).

This example creates several endpoints which we can hit in a browser.

Try hitting [`http://localhost:8080/`](http://localhost:8080/) and you'll get this back:

```json
{
    "links": {
        "book index": "http://localhost:8080/book",
        "root": "http://localhost:8080/",
        "self": "http://localhost:8080/"
    }
}
```

This shows us all of the top-level endpoints that are available.

Now if you hit the `book index` url [`http://localhost:8080/book`](http://localhost:8080/book) you'll see a list of
book objects:

```json
[
    {
        "Author": "Charles Dickens",
        "ISBN": "345678",
        "Id": 0,
        "Name": "Great Expectations",
        "links": {
            "book index": "http://localhost:8080/book",
            "root": "http://localhost:8080/",
            "self": "http://localhost:8080/book/0"
        }
    },
    {
        "Author": "Daniel Dafoe",
        "ISBN": "234567",
        "Id": 1,
        "Name": "Robinson Crusoe",
        "links": {
            "book index": "http://localhost:8080/book",
            "root": "http://localhost:8080/",
            "self": "http://localhost:8080/book/1"
        }
    },
    {
        "Author": "William Shakespeare",
        "ISBN": "123456",
        "Id": 2,
        "Name": "Henry V",
        "links": {
            "book index": "http://localhost:8080/book",
            "root": "http://localhost:8080/",
            "self": "http://localhost:8080/book/2"
        }
    }
]
```

Each of the objects has a link for:

- themselves
- the index for their type
- a link to the root of the server

These links make the APIs very easy to navigate.

In addition to the GET endpoints described above, burrow also exposes PUT endpoints for updating objects:

```bash
$ curl -X PUT -d '{"Name": "Not-So-Great Expectations"}' http://localhost:8080/book/0
```

We get back the JSON object for the record we just updated:

## More Examples

Burrow has the ability to create links to related objects. See
[examples/references.go](http://github.com/zfjagann/burrow/tree/master/examples/references.go).

## Dependencies

Burrow is based on [Traffic](http://github.com/pilu/traffic), and requires that it be installed to be used.
