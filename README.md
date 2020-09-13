
![Algolia logo](https://upload.wikimedia.org/wikipedia/commons/thumb/6/69/Algolia-logo.svg/1280px-Algolia-logo.svg.png)


# GoPlaces

Simple wrapper for Algolia Places API with extended functionality

## Features

- No need to manually decode API responses
- Addresses extracting from hits
- Generating label from address object, f.e. "Brickell Avenue, Miami, Florida, 33131"
- 0 deps

## Usage

```go
package main

import (
    "fmt"
    "github.com/yuriizinets/goplaces"
)

func main() {
    // (Optional) Set AppID and AppKey for auth.
    // Also you can provide that values inside parameters for per-request credentials usage.
    // Free tier does not require that at all, so credentials are optional in both cases
    goplaces.AppID = "..."
    goplaces.AppKey = "..."
    // Resp is QueryResponse object, that represents Algolia API response as-is
    resp, err := goplaces.Query(goplaces.Parameters{
        Query: "Brickell Avenue, Miami, Florida",
    })
    if err != nil {
        panic(err)
    }
    // Raw response is not so useful and uncomfortable to use
    // You can extract addresses (Address objects) from Hit records for comfortable work
    addresses := goplaces.ExtractAddresses(resp.Hits)
    fmt.Println(addresses)
    // Also, you can create label from address
    // It will drop missing values from address object
    label := NewLabelFromAddress(addresses[0])
}
```
