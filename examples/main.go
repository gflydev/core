package main

import (
	"fmt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
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
	// http://localhost:7789/api/v1/info?param1=one&arr[]=item1&arr[]=item2&arr[]=3
	queryData := make(core.Data)
	if err := c.ParseQuery(&queryData); err != nil {
		return err
	}

	keys := queryData.Keys()
	for _, key := range keys {
		value := queryData.Get(key)

		switch value.(type) {
		case string:
			log.Infof("%s: %s", key, value.(string))
			break
		case []string:
			for _, v := range value.([]string) {
				log.Infof("%s[]: %s", key, v)
			}
			break
		}
	}

	response := make(core.Data).
		Set("name", core.AppName).
		Set("server", core.AppURL).
		Set("query", queryData)

	return c.JSON(response)
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
	return c.View("home", core.Data{
		"title": "gFly | Laravel inspired web framework written in Go",
	})
}

// =========================================================================================
//                                     Routers
// =========================================================================================

func router(g core.IFly) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", ""),
		utils.Getenv("API_VERSION", ""),
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

	// Register router
	app.RegisterRouter(router)

	app.Run()
}
