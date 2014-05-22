Burrow
======

Golang Restful API Framework
=======
# Burrow - Golang Restful API Framework

Burrow is a REST API Framework designed with:

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

- Automatically generate API endpoints
- Auto-generate documentation for your API endpoints *(coming soon)*
- Manage links between referenced objects
- Generate links in your API objects to make navigating your API dead simple

## Introduction

Consider the following type:

    type Book struct {
        Id     int `rest:"id"`
        Name   string
        ISBN   string
        Author string
    }

In addition to the type itself, Burrow needs to know how to retrieve objects of this type.
To do this, we will define two functions.

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

The first function returns a book based on its id. The second gets a list of all the books.

In a more realistic application, these would likely be DB accessors, but for the sake of example, they'll just read
from a global variable `AllBooks`, which we will initialize in the `init` function:

    func init() {
        AllBooks = make([]Book, 3)
        AllBooks[0] = Book{0, "Great Expectations", "345678", "Charles Dickens"}
        AllBooks[1] = Book{1, "Robinson Crusoe", "234567", "Daniel Dafoe"}
        AllBooks[2] = Book{2, "Henry V", "123456", "William Shakespeare"}
    }

Finally, we create a Burrow api and use it:

    func main() {
        api := burrow.NewApi()
        api.AddApi(Book{}, GetBook, GetBooks)
        api.Serve("0.0.0.0", 8080)
    }

This will create a couple API endpoints that we can hit.

The full code can be found [here](http://github.com/zfjagann/burrow/tree/master/examples/simple.go).

Since Burrow is managing all the API endpoints for us, it can give us a nice index if you hit the root of the webserver.

So if you hit [`http://localhost:8080/`](http://localhost:8080/) you'll see this:

    {
        "links": {
            "book index": "http://localhost:8080/book",
            "root": "http://localhost:8080/",
            "self": "http://localhost:8080/"
        }
    }

This shows us all of the top-level endpoints that are available.

Now if you hit the `book index` url [`http://localhost:8080/book`](http://localhost:8080/book) you'll see a list of 
all of our book objects:

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

Each of the objects has a link for:

- themselves
- the index for their type
- a link to the root of the server

These links make the APIs very easy to navigate.

Burrow also has the ability to create links to related objects. For an example of this, see
[examples/references.go](http://github.com/zfjagann/burrow/tree/master/examples/references.go).

## Dependencies

Burrow is based on [Traffic](http://github.com/pilu/traffic), and requires that it be installed to be used.
