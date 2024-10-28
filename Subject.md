# Coding Challenge Guidelines

Please organize, design, test, document and write your code as if it were going into production, then send us a link to the hosted repository (e.g. Github, Bitbucket...).

## Functional spec

The goal of this exervice is to create an HTTP API server that is able
to produce aggregations based on the data rendered by the Upfluence SSE API.

Your API server should listen on the port `8080` on your local loopback.
And accept **ONLY** HTTP `GET` request for the path `/analysis`.

All other HTTP requests should result as a `404` response.

Your API should take as query parameters two values:

* `duration`: as a unit of time (5s, 10m, 24h are all valid input) for
  which your API should analyze the post outputed by our SSE API.
* `dimension`: The value we want to generate the statistics upon, it
  could any of the following:
    * `likes`
    * `comments`
    * `favorites`
    * `retweets`

Your API should return a JSON payload with the following informations:

* the total number of posts analyzed
* the minimum timestamp of the posts gathered during the analysis
* the maximum timestamp of the posts gathered during the analysis
* the average value of a dimension from the posts

### Input stream

Upfluence has a publicly available HTTP API endpoint streaming some of
the posts processed by our system in real-time (It uses
[Server Sent Events (SSE)](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events)).

You can find this stream here: https://stream.upfluence.co/stream

#### Stream payload description

Each "event" sent by  the stream endpoint is a `JSON` encoded object that represents a single post being detected by the upfluence platform. The payload  looks like the following for an instagram media.

```
{
   "instagram_media": {
      "id": 102810280182,
      "text": "Some text",
      "likes": 27,
      "comments": 42,
      // .... some other fields
   }
}
```

The (single) root key is the type of social post, It can either be:

* `pin`
* `instagram_media`
* `youtube_video`
* `article`
* `tweet`
* `facebook_status`

The value is a JSON object representing the values with some details
about the post. Its structure can vary from one social media post to
another. That being said all values include the key `timestamp`
representing the creation date of the post (encoded as a UNIX timestamp).

You can easily visualize the stream content with the following command
from your terminal:

```bash
 curl 'https://stream.upfluence.co/stream'
```

### Example of API usage

```
$ curl localhost:8080/analysis?duration=30s&dimension=likes
// Should block for 30 seconds

200 OK
Content-Type: application/json
Content-Length: ...
// Some other headers

{
  "total_posts": 20
  "minimum_timestamp": X,
  "maximum_timestamp": Y,
  "avg_likes": 50
}
```

## Expected ouput

### Readme

Write your README as if it was for a production service. Include the following items:

* Description of your solution.
* Reasoning behind your technical choices, including architectural.
* Trade-offs you might have made, anything you left out, or what you might do differently if you were to spend additional time on the project.
* And of course, how to run your project. Including expected dependencies and command examples showing us how to install and run your project.

### Technical solution

In order to complete this challenge you can use the language you want.

You are only allowed to use the standard library (stdlib) of the
language in order to complete the assignment except for the HTTP server
/ request handling (i.e. you can use flask or django if you use python or rails / sinatra in ruby).

### How we review

Your application will be reviewed by at least three of our engineers. We do take into consideration your experience level.

We value quality over feature-completeness. It is fine to leave things aside provided you call them out in your project's `README`. The goal of this code sample is to help us identify what you consider production-ready code. You should consider this code ready for final review with your colleague, i.e. this would be the last step before deploying to production.

The aspects of your code we will assess include:

* **Architecture**: How clean is the separation between the different components? How well are the different classes separated?
* **Clarity**: Does the README clearly and concisely explain the problem and solution? Are technical trade-offs explained?
* **Correctness**: Does the application do what was asked? If there is anything missing, does the README explain why it is missing?
* **Code quality**: Is the code simple, easy to understand, and maintainable? Are there any code smells or other red flags? Is the coding style consistent with the language's guidelines? Is it consistent throughout the codebase?
* **Technical choices**: do choices of libraries, architecture etc. seem appropriate for the chosen application?

Bonus point (those items are optional):

* **Testing quality**: How extensive is your testing coverage? Is your testing approach pragmatic?
* **Efficiency**: What is the footprint of your program? Both in terms of CPU and memory usage
* **Scalability**: Will technical choices scale well (increase both in terms of post throughput and number of concurrent client) ? If not, is there a discussion of those choices in the README?