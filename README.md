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


