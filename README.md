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

The following go code will define an API for managing books, complete with links and endpoints to update records:

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

func GetBook(id int) (interface{}, error) {
    if id < 0 || id >= len(AllBooks) {
        return nil, burrow.ApiError(404, "Could not find book with id: ", id)
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

func init() {
    // Create dummy data...
    AllBooks = make([]Book, 3)
    AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens"}
    AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe"}
    AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare"}
}

func main() {
    api := burrow.NewApi()

    api.Add(burrow.New(Book{}, nil, GetBook, GetBooks, UpdateBook, nil))

    // Run the server!
    api.Serve("0.0.0.0", 8080)
}
```

For the complete example see
[examples/simple.go](http://github.com/zfjagann/burrow/blob/develop/examples/simple.go).

### The API

We can start using the API immediately. Start by hitting the root of the server: `http://localhost:8080/`.

The response JSON looks like:

```json
{
    "links": {
        "Book index": "http://localhost:8080/book",
        "root": "/",
        "self": "/"
    }
}
```

The links provided in the output indicate which endpoints are available.

Now follow the `Book index` link:

```json
[
    {
        "Author": "Charles Dickens",
        "ISBN": "345678",
        "Id": 0,
        "Name": "Great Expectations",
        "links": {
            "Book index": "http://localhost:8080/book",
            "root": "/",
            "self": "http://localhost:8080/book/0"
        }
    },
    {
        "Author": "Daniel Dafoe",
        "ISBN": "234567",
        "Id": 1,
        "Name": "Robinson Crusoe",
        "links": {
            "Book index": "http://localhost:8080/book",
            "root": "/",
            "self": "http://localhost:8080/book/1"
        }
    },
    {
        "Author": "William Shakespeare",
        "ISBN": "123456",
        "Id": 2,
        "Name": "Henry V",
        "links": {
            "Book index": "http://localhost:8080/book",
            "root": "/",
            "self": "http://localhost:8080/book/2"
        }
    }
]
```

We get back our dummy book data that we setup. Each of the records also has a link to themselves to allow you to show an individual record.

Since we defined `UpdateBook` and handed it to the API, we can also update book objects:

```bash
$ curl -XPUT -d '{"Author":"Willem Dafoe"}' http://localhost:8080/book/1
```

Gives us back:
```json
{
    "Author": "Willem Dafoe",
    "ISBN": "234567",
    "Id": 1,
    "Name": "Robinson Crusoe",
    "links": {
        "Book index": "http://localhost:8080/book",
        "root": "/",
        "self": "http://localhost:8080/book/1"
    }
}
```


### The CRUD

The CRUD model defined above for managing books includes functions for viewing and modifying books. This is
sufficient for a simple application, but a more thorough CRUD model would also allow deletion and addition of books.

The `burrow.CRUD` interface provides an interface for users to define full CRUDs.

The `CRUD` includes 7 methods that can be used to manipulate objects.

```go
type CRUD interface {
    /*
        Returns the name of the type that this CRUD represents.

        This is used to determine URIs.
    */
    Name() string

    /*
        Returns the reflected type of the object this CRUD manages.
    */
    Reflect() reflect.Type

    /*
        Create a new object and return it.
    */
    Create() (interface{}, error)

    /*
        Update a given object.
        Returns nil if successful, and an error otherwise.
    */
    Update(interface{}) error

    /*
        Load and return a single object.
        The error return value will be nil if successful, and non-nil otherwise.
    */
    Read(int) (interface{}, error)

    /*
        Load and return all objects.
        The error return value will be nil if succesful, and non-nil otherwise.
    */
    Index() ([]interface{}, error)

    /*
        Delete the object specified by the given id.
        Returns nil if successful, and an error otherwise.
    */
    Delete(int) error
}
```


## More Examples

Burrow has the ability to create links to related objects. See
[examples/references.go](http://github.com/zfjagann/burrow/tree/master/examples/references.go).

## Dependencies

Burrow is based on [Traffic](http://github.com/pilu/traffic), and requires that it be installed to be used.
