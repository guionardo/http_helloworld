# http_helloworld

Minimal HTTP server for testing purposes on docker

[![CodeQL](https://github.com/guionardo/http_helloworld/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/guionardo/http_helloworld/actions/workflows/codeql-analysis.yml)
[![CI](https://github.com/guionardo/http_helloworld/actions/workflows/ci.yaml/badge.svg)](https://github.com/guionardo/http_helloworld/actions/workflows/ci.yaml)
[![Go](https://github.com/guionardo/http_helloworld/actions/workflows/go.yml/badge.svg)](https://github.com/guionardo/http_helloworld/actions/workflows/go.yml)

This server exposes only a root endpoint, with a JSON response like above: 

```json
{
  "ip": "172.19.0.1:45370",
  "runningTime": "10.563721891s",
  "startTime": "2022-03-12 11:47:57.966620544 +0000 UTC m=+0.000435499",
  "time": "2022-03-12 11:48:08.530328605 +0000 UTC m=+10.564143553",
  "requestCount": 3,
  "tag": "testing"
}
```

## Usage

```bash
docker run --rm -p 8080:8080 guionardo/http_helloworld:latest
```

This server will listen to HTTP port 8080, but you can change this behavior using a environenment variable PORT.

```bash
docker run --rm p 8081:8081 -e PORT=8081 -e TAG=testing guionardo/http_helloworld:latest
```
You can use the environment vars:

- PORT = Por number (between 1 and 65535)
- TAG = Text for add into JSON response, use for show any custom data


## Custom endpoints

You can add your customized static endpoints using a configuration file like this:

In folder _./custom_responses_ add a _routes.json_ file with an array of objects:

| Field        | type    | Description                                               |
|--------------|---------|-----------------------------------------------------------|
| path         | string  | URL path                                                  |
| source_file  | string  | name of file with content, without path                   |
| method       | string  | HTTP method (default = GET)                               |
| status_code  | integer | HTTP status code (default = 200)                          |
| content_type | string  | Content-Type header (default = detected from source file) |

Example:

```json
[
    {
        "path": "/api",
        "source_file": "api.content",
        "method": "GET",
        "status_code": 200,
        "content_type": "application/json",
    },
    {
      "path":"/test",
      "source_file": "test.content",
      "method": "POST",
      "status_code": 202
    }
]
```

And add the files _api.content_ and _test.content_ to _./custom_responses_ folder with some content.

When the server receives a _/api_ request, it will response the body of _api.content_ file. And so on to _/test_ : _test.content_

You need to change the docker command to this:

```bash
docker run --rm -p 8080:8080 -e PORT=8080 \
  -v ${PWD}/custom_responses:/app/custom_responses \
  guionardo/http_helloworld:latest
```