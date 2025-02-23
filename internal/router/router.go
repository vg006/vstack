//go:build !js || !wasm

package types

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/fsnotify.v1"
)

type Router struct {
	Dir       string
	PagesDir  string
	PublicDir string
	Routes    []*Route
	Page500   string
}

type Project struct {
	Name    string
	ModName string
	SrcDir  string
}

type Route struct {
	template *template.Template

	UrlPath  string
	FilePath string
	Page200  string
}

func NewRouter(dir string) *Router {
	return &Router{
		Dir:       dir,
		PagesDir:  dir + "/pages",
		PublicDir: dir + "/public",
	}
}

func (r *Route) Render(data ...any) error {
	var buf bytes.Buffer
	err := r.template.Execute(&buf, data)
	if err != nil {
		return err
	}
	r.Page200 = buf.String()
	return nil
}

func (router *Router) Load() error {
	var routes []*Route

	err := filepath.WalkDir(router.PagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Base(path) != "200.html" {
			return nil
		}

		route := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(router.PagesDir))
		if strings.HasPrefix(route, "/") {
			route = strings.TrimPrefix(route, "/")
		}
		route = strings.TrimSuffix(route, "200.html")
		route = strings.TrimSuffix(route, "/")
		if route == "" {
			route = "/"
		} else if !strings.HasPrefix(route, "/") {
			route = "/" + route
		}

		tmpl, err := template.ParseFiles(path)
		if err != nil {
			return err
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, nil)
		if err != nil {
			return err
		}

		routes = append(routes, &Route{
			template: tmpl,
			UrlPath:  route,
			FilePath: path,
			Page200:  buf.String(),
		})
		return nil
	})

	if err != nil {
		return err
	}
	router.Routes = routes

	page500Path := filepath.Join(router.PagesDir, "500.html")
	if _, err := os.Stat(page500Path); os.IsNotExist(err) {
		baseHTML := `<html><head><title>500 - Internal Server Error</title></head><body><h1>500 - Internal Server Error</h1><p>Something went wrong.</p></body></html>`
		router.Page500 = baseHTML
	} else {
		tmpl, err := template.ParseFiles(page500Path)
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		err = tmpl.Execute(&buf, nil)
		if err != nil {
			return err
		}
		router.Page500 = buf.String()
		defer buf.Reset()
	}
	return nil
}

func (r *Router) Reload() error {
	for _, route := range r.Routes {
		err := route.Reload()
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Route) Reload() error {
	tmpl, err := template.ParseFiles(r.FilePath)
	if err != nil {
		return err
	}
	r.template = tmpl
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		return err
	}
	r.Page200 = buf.String()
	defer buf.Reset()
	return nil
}

func (r *Router) InitHMR() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating watcher: %v\n", err)
		return err
	}

	err = filepath.WalkDir(r.PagesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return fmt.Errorf("error watching %s: %v", path, err)
			}
			fmt.Fprintf(os.Stdout, "Watching directory: %s\n", path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error setting up watcher: %v", err)
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".html") {
					err := r.Reload()
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error reloading templates: %v\n", err)
						return
					}
					fmt.Fprintf(os.Stdout, "Reloaded templates due to change in %s\n", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Fprintf(os.Stderr, "Watcher error: %v\n", err)
			}
		}
	}()

	return nil
}
