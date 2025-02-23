//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	vstack "github.com/vg006/vstack/pkg"
)

func main() {
	counter := 0

	c := vstack.NewClient()
	if err := c.AddIslands(vstack.Island{
		Id: "add",
		Handlers: map[string]js.Func{
			"click": js.FuncOf(func(this js.Value, args []js.Value) any {
				counter++
				go func() {
					url := "http://localhost:8080/api/todo"
					data := make([]byte, 0, 64)
					data = fmt.Appendf(data, `{"task": "%d"}`, counter)
					req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
					if err != nil {
						fmt.Println("Request error:", err)
						return
					}
					req.Header.Set("Content-Type", "application/json")

					resp, err := http.DefaultClient.Do(req)
					if err != nil {
						fmt.Println("POST error:", err)
						return
					}
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						fmt.Println("Read error:", err)
						return
					}

					todoList := c.GetElementById("todoList")
					newItem := c.CreateElement("div")
					newItem.SetInnerHTML(string(body))
					todoList.AppendChild(newItem)
				}()
				err := c.Reload()
				if err != nil {
					return err
				}
				return nil
			}),
		},
	}, vstack.Island{
		Class: "delete",
		Handlers: map[string]js.Func{
			"click": js.FuncOf(func(this js.Value, args []js.Value) any {
				t := args[0].Get("target")
				for t.Type() == js.TypeObject && t.Get("tagName").String() != "DIV" {
					t = t.Get("parentElement")
					if t.IsNull() {
						break
					}
				}
				if !t.IsNull() && t.Type() != js.TypeUndefined {
					t.Get("parentNode").Call("removeChild", t)
				}
				err := c.Reload()
				if err != nil {
					return err
				}
				return nil
			}),
		},
	}); err != nil {
		fmt.Println(err.Error())
	}
	select {}
}
