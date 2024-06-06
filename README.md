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

| Name    | HTTP Method | Route          |
|---------|-------------|----------------|
| Health  | GET         | /healthy       |
| Index   | GET         | /              |
| Summary | POST        | /summary       |

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

* [docker](https://www.docker.com)
* [Go](https://golang.org/) (optional, required in case of building from source)

### Installing

1. Clone this repository to your local machine.
```
git clone https://github.com/rohit4813/web-analyser.git
```
2. Install dependencies
```
go get -v ./...
```
3. Start the application via docker:
```
docker-compose up
```


### Building From Source
1. Clone this repository to your local machine.
```
git clone https://github.com/rohit4813/web-analyser.git
```
2. Install dependencies
```
go get -v ./...
```
3. Run the golang application:
```
go build cmd/web/main.go && ./main
```

## Running the tests


## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Assumptions
* I have assumed that we do not have to select the package that itself is the most performant, 
my main priority was to solve the problem most efficiently with the packages I have chosen
* Taking w3 org for considering standard ways to define html version: https://www.w3.org/QA/2002/04/valid-dtd-list.html
* For deciding the number of external links, internal links and inaccessible links, I am considering only anchor links, 
 and assuming internal links as links having the host empty or same as that of the given url, external links as links
 having different host, inaccessible links as links which are not in the proper format as per the url package. I am
 not checking whether the inaccessible links are reachable over the internet by trying to get/dial the url.
* The user data has to be sent via POST method to the backend.
* In case of unreachable link, it is possible that the user internet connectivity is poor or the entered url dns 
  resolution fails, in this case we will not have the http status code, so I am just displaying the error message.

## Improvements
* We can use goroutines, channels and wait groups to analyse different aspects of the URL concurrently.
* We can try to get/dial the links to add it to the inaccessible links map, though care must be taken 
  on how many links and how much timeout for each link should be given since we might have to end up 
  waiting to populate it.
* URL regex is basic, it can be improved to incorporate more patterns
* We can use godoc to write documentation more effectively.
* Templates can be tested for different screen sizes.
* Errors can be displayed in the index page itself for smooth user experience.
