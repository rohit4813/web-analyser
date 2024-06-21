# Web Analyzer

## Introduction

The objective was to build a web application that does an analysis of a web-page/URL.


## Features

The application show a form with a text field in which a user can type in the URL of the webpage to be analysed.
After analysing the URL, the user is presented with the summary of the html page which includes the below fields:

* HTML Version
* Title
* Headers count of each level
* Internal Links count
* External Links count
* Inaccessible Links count
* Has login form

## Endpoints

| Name                | HTTP Method | Route          |
|---------------------|-------------|----------------|
| Index Page          | GET         | /              |
| Summary Page        | POST        | /summary       |
| Health Api Endpoint | GET         | /healthy       |

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* [docker](https://www.docker.com)
* [Go](https://golang.org/) (optional, required in case of building from source)

### Usage

1. Clone this repository to your local machine.
```
git clone https://github.com/rohit4813/web-analyser.git
```
2. Start the application via docker:
```
docker-compose up
```
3. Visit http://localhost:8080/ in any browser.

### Building From Source
1. Clone this repository to your local machine.
```
git clone https://github.com/rohit4813/web-analyser.git
```
2. Install dependencies:
```
go get -v ./...
```
3. Run the golang application:
```
go build cmd/web/main.go && ./main
```
4. Visit http://localhost:8080/ in any browser.


## Running the tests
```
go test -v ./...
```

## Project structure

```shell
web-analyser
├── api
│  ├── backend
│  │  ├── analyser
│  │  │  ├── analyser.go
│  │  │  ├── analyser_test.go
│  │  │  ├── handler.go
│  │  │  ├── handler_test.go
│  │  │  ├── model.go
│  │  │  └── template.go
│  │  └── health
│  │     └── health.go
│  ├── router
│  │  ├── middleware
│  │  │  ├── request_id.go
│  │  │  ├── request_id_test.go
│  │  │  └── request_log.go
│  │  └── router.go
│  └── templates
│     ├── error.gohtml
│     ├── index.gohtml
│     └── summary.gohtml
├── cmd
│  └── web
│     └── main.go
├── config
│  └── config.go
├── internal
│  └── utils
│     ├── ctx
│     │  ├── ctx.go
│     │  └── ctx_test.go
│     ├── error
│     │  └── error.go
│     ├── html
│     │  ├── html.go
│     │  └── html_test.go
│     ├── http
│     │  ├── http.go
│     │  └── http_test.go
│     └── logger
│        └── logger.go
├── mocks
│  ├── analyser_mock.go
│  ├── http_mock.go
│  └── template_mock.go
├── .dockerignore
├── .env
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Assumptions
* I have assumed that we do not have to select the package that itself is the most performant, 
my main priority was to solve the problem most efficiently with the packages I have chosen
* Taking w3 org for considering standard way to define html version: https://www.w3.org/QA/2002/04/valid-dtd-list.html
* For deciding the number of internal links, external links and inaccessible links, I am considering only anchor links
 without mailto, tel, javascript and assuming internal links as links having the host empty or same as that of the 
 given url, external links as links having different host, inaccessible links as links which are not in the proper 
 format as per the url package. I am not checking whether the inaccessible links are reachable over the internet by 
 trying to get/dial the url.
* The user data has to be sent via POST method to the backend.
* There may be ways to render the templates more effectively rather than loading all the templates in memory beforehand.
  I assumed the template rendering performance was not critical for this task.
* In case of unreachable link, it is possible that the user internet connectivity is poor or the entered url dns
  resolution fails, in this case we will not have the http status code, so I am just displaying the error message.

## Improvements
* We can use goroutines, channels and wait groups to analyse different aspects of the URL concurrently.
* We can try to get/dial the links to add it to the inaccessible links map, though care must be taken 
  on how many links and how much timeout for each link should be given since we might have to end up 
  waiting to populate it. There needs to be a maximum limit on the number of get/dial we should
  do as well.
* URL regex is basic, it can be improved to incorporate more patterns.
* More documentation can be added.
* Errors can be displayed in the index page itself for smooth user experience.

