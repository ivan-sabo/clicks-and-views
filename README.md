# Views and Clicks API

## About project

This is simple Golang RESTful API project.  
It consists of two resources, Click and View.

To start a project simply navigate to "cmd" folder and run

```console
foo@bar:~$ go run main.go
```

... or from root folder:

```console
foo@bar:~$ make run
```

To run local OpenAPI web client, navigate to project root folder run the following  
command, and open http://localhost in your browser:

```console
foo@bar:~$ make swagger
```

## Additional documentation

Besides OpenAPI documentation, there is only one source of documentation:  
the code itself :)

This project follows all guidelines of https://go.dev/doc/effective_go and  
https://go.dev/wiki/CodeReviewComments.

All public methods and properties contain doc comments.

## TODO

Some important features are not implemented:

1. Logging and Error handling (user-friendly error messages)
2. Pagination
3. Sorting
