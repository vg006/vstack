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

func (r *Route) Render(data any) (string, error) {
	var buf bytes.Buffer
	err := r.template.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (r *Route) Reload() error {
	tmpl, err := template.ParseFiles(r.FilePath)
	if err != nil {
		return err
	}
	r.template = tmpl
	return nil
}

func Reload(routes []Route) error {
	for _, route := range routes {
		if err := route.Reload(); err != nil {
			return err
		}
	}
	return nil
}

func Load(dir string, html string, hmr bool) ([]Route, error) {
	var routes []Route
	var watcher fsnotify.Watcher

	if hmr {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return nil, fmt.Errorf("error creating watcher: %v", err)
		}
		defer watcher.Close()

		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".html") {
						Load(dir, html, hmr)
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

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		if hmr && d.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return fmt.Errorf("error watching %s: %v", path, err)
			}
			fmt.Fprintf(os.Stdout, "Watching directory: %s\n", path)
		}

		route := strings.TrimPrefix(filepath.ToSlash(path), filepath.ToSlash(dir))
		if strings.HasPrefix(route, "/") {
			route = strings.TrimPrefix(route, "/")
		}
		route = strings.TrimSuffix(route, html)
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
		return nil, err
	}
	return routes, nil
}

func InitiateHMR(dir, html string, hmr bool) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %v", err)
	}
	defer watcher.Close()

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
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
		return fmt.Errorf("error setting up watcher: %v", err)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write && strings.HasSuffix(event.Name, ".html") {
					Load(dir, html, hmr)
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
