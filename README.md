# gFly Core - Core Web written in Go

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

### Setup `gFly Core`
```bash
mkdir myweb && cd myweb
go mod init myweb
go get -u github.com/gflydev/core@latest
```

### Play with gFly

#### Create folder 

Inside your application folder `myweb` 
```bash
mkdir -p storage/tmp
mkdir -p storage/logs
mkdir -p storage/app
mkdir -p resources/views
mkdir -p resources/tls
```

Create static page
```bash
mkdir public
touch public/index.html
```
Content `index.html`
```html
<!DOCTYPE html>
<html lang="en" dir="ltr">
<head>
    <meta charset="UTF-8"/>
    <meta http-equiv="X-UA-Compatible" content="IE=edge"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no"/>
    <title>gFly | Laravel inspired web framework written in Go</title>
</head>
<body>
    <h2>gFly | Laravel inspired web framework written in Go</h2>
</body>
</html>
```

#### Create app `main.go`
```go
package main

import (
    "fmt"
    "github.com/gflydev/core"
    "github.com/gflydev/core/utils"
    _ "github.com/joho/godotenv/autoload"
)

// =========================================================================================
//                                     Default API
// =========================================================================================

// NewDefaultApi As a constructor to create new API.
func NewDefaultApi() *DefaultApi {
    return &DefaultApi{}
}

// DefaultApi API struct.
type DefaultApi struct {
    core.Api
}

func (h *DefaultApi) Handle(c *core.Ctx) error {
    return c.JSON(core.Data{
        "name":   core.AppName,
        "server": core.AppURL,
    })
}

// =========================================================================================
//                                     Home page 
// =========================================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
    return &HomePage{}
}

type HomePage struct {
    core.Page
}

func (m *HomePage) Handle(c *core.Ctx) error {
    return c.HTML("<h2>Hello world</h2>")
}

// =========================================================================================
//                                     Routers
// =========================================================================================

func router(g core.IFlyRouter) {
    prefixAPI := fmt.Sprintf(
        "/%s/%s",
        utils.Getenv("API_PREFIX", "api"),
        utils.Getenv("API_VERSION", "v1"),
    )

    // API Routers
    g.Group(prefixAPI, func(apiRouter *core.Group) {
        apiRouter.GET("/info", NewDefaultApi())
    })

	// Web Routers
    g.GET("/home", NewHomePage())
}

// =========================================================================================
//                                     Application 
// =========================================================================================

func main() {
    app := core.New()

    // Register middleware
    //app.RegisterMiddleware(hookMiddlewares)

    // Register router
    app.RegisterRouter(router)

    app.Run()
}
```

#### Run and Check

Run
```bash
go run main.go
```

Check API
```bash
curl -X 'GET' \
  'http://localhost:7789/api/v1/info' | jq
```

Note: Install [jq](https://jqlang.github.io/jq/) tool to view JSON format

Check static page
http://localhost:7789/index.html

Check dynamic page
http://localhost:7789/home