# GOHTTP

GOHTTP is a lean and highly configurable Go HTTP client library with advanced features like rate limiting, circuit breaker,
retry mechanisms and endpoint-specific configurations. Only use the features you need.

## Usage
The library uses a `Request` struct which is a wrapper around `http.Request` type. The `Request` type can be used with
extra configuration options following the `RequestConfig` type.

A sample `GET` request would look like the following.

```go
package main

import (
    "log"
    "net/http"
)

const url = "https://jsonbox.io/demobox_6d9e326c183fde7b"

func main() {

    client := gohttp.NewHtpClient(http.DefaultClient, gohttp.WithTimeout(30*time.Second))    

    // we are going to tell the api that we accept json as valid response
    options := []gohttp.RequestConfig{gohttp.WithAcceptJSONHeader()}
    response, err := client.Get(url, options...)
    if err != nil {
        log.Println(err)
        return
    }
    log.Println(response)
}
```
