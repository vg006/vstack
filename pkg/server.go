//go:build !js || !wasm

package vstack

import (
	"github.com/vg006/vstack/internal/server"
	utils "github.com/vg006/vstack/internal/utils"
)

type (
	Router = server.Router
	Route  = server.Route
)

func NewRouter(dir string) *Router {
	return &Router{
		Dir:       dir,
		PagesDir:  dir + "/pages",
		PublicDir: dir + "/public",
	}
}

func EchoPath(path string) string {
	return utils.EchoPath(path)
}

func FiberPath(path string) string {
	return utils.FiberPath(path)
}
