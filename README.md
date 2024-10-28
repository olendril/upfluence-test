# Upfluence Test

## The Solution 

The API server listen on the port `8080` and accept **ONLY** HTTP `GET` request for the path `/analysis`.

The API should take as query parameters two values:

* `duration`: as a unit of time (5s, 10m, 24h are all valid input) represent the duration for which the API listen for 
the upfluence API
* `dimension`: The value we want to generate the statistics upon, it
  could any of the following:
    * `likes`
    * `comments`
    * `favorites`
    * `retweets`

The API return a JSON payload with the following informations:

* the total number of posts analyzed
* the minimum timestamp of the posts gathered during the analysis
* the maximum timestamp of the posts gathered during the analysis
* the average value of a dimension from the posts

## TODO

- I need to set up unit test and integration test.
- I want to implement a clean infrastructure for the project by implementing a standard architecture for go project
https://github.com/golang-standards/project-layout
- We need to rework the media parsing from the upfluence API to have one structure by media type to allow finer parsing 
and have a matching of similar field (like retweet and share) 
- I want to set up an open api contract with a generator that generate the server side code to allow a better maintainability and facilitate the usability of the service.

## How to Run

The easy option to run this project is to launch using docker compose.

````shell
docker compose up
````

You can also launch the project using the go cli.

````shell
go run main.go
````


