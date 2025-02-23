//go:build !js || !wasm

// * NOTE! The above is not a comment, it is a go:directive that tells the compiler about the type of supported architecture upon building
package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	vstack "github.com/vg006/vstack/pkg"
)

func main() {
	// Initiates a echo server instance
	e := echo.New()
	e.Use(middleware.CORS())

	// Intiates of vstack router
	routes := vstack.NewRouter("src")

	// Initiates Watcher
	err := routes.InitHMR()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Fetches all the routes and their associated pages(200.html)
	err = routes.Load()
	if err != nil {
		fmt.Println(err.Error())
	}

	// Custom error handler
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}
		// Sends the 500.html in the 'src/pages' dir, otherwise returns the default(hardcoded) html
		c.HTML(500, routes.Page500)
	}

	// Serves the public assets
	e.Static("/public/", routes.PublicDir)

	// Registering the routes for the webpages
	for _, route := range routes.Routes {
		// Returns the echo-specific route
		path := vstack.EchoPath(route.UrlPath)
		e.GET(path, func(c echo.Context) error {
			// Handling multipe routes
			switch path {
			case "/home/:id":
				err := route.Render(c.Param("id"))
				if err != nil {
					fmt.Println(err.Error())
					return err
				}
			}
			// Sends the 200.html in the route's dir, otherwise returns the error page
			return c.HTML(200, route.Page200)
		})
	}

	// Registers a dummy api endpoint
	e.POST("/api/todo", func(c echo.Context) error {
		var todo map[string]string
		if err := c.Bind(&todo); err != nil {
			return c.HTML(400, fmt.Sprintf("<h1>Error</h1>"))
		}

		fmt.Println("Added todo:", todo["task"])
		return c.HTML(200, fmt.Sprintf("<div class=\"flex flex-row flex-nowrap justify-between p-3\"><li class=\"text-lg font-medium text-gray-700\">%s</li><button class=\"delete bg-red-500 text-white px-4 py-2 rounded\">Delete</button></div>", todo["task"]))
	})

	// Obviously, Starts the server
	e.Logger.Fatal(e.Start(fmt.Sprintf(":8080")))
}
