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
	Dir    string
	Html   string
	Routes []Route
	Signal chan int
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
	Markup   string
}

func NewRouter(dir, html string) *Router {
	return &Router{
		Dir:    dir,
		Html:   html,
		Signal: make(chan int, 1),
	}
}

func (r *Route) Render(data any) (string, error) {
	var buf bytes.Buffer
	err := r.template.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (router *Router) Load() error {
	var routes []Route

	err := filepath.WalkDir(router.Dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		route := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(router.Dir))
		if strings.HasPrefix(route, "/") {
			route = strings.TrimPrefix(route, "/")
		}
		route = strings.TrimSuffix(route, router.Html)
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

		routes = append(routes, Route{
			template: tmpl,
			UrlPath:  route,
			FilePath: path,
			Markup:   buf.String(),
		})
		return nil
	})

	if err != nil {
		return err
	}
	router.Routes = routes
	return nil
}

func (r *Router) InitHMR() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating watcher: %v\n", err)
		return
	}

	err = filepath.WalkDir(r.Dir, func(path string, d fs.DirEntry, err error) error {
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
		return
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".html") {
					err := r.Load()
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
}
