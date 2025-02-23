//go:build js && wasm

package client

import (
	"fmt"
	"syscall/js"
)

type Client struct {
	Dom     js.Value
	Islands []Island
}

type Node struct {
	Js js.Value
}

type Island struct {
	Value      any
	DataSource string
	Handlers   map[string]js.Func
	Id         string
	Class      string
}

func (c *Client) AddIslands(is ...Island) error {
	for _, i := range is {
		c.Islands = append(c.Islands, i)

		if i.Id == "" && i.Class == "" {
			return fmt.Errorf("Element id or class not found")
		}

		if i.Id != "" {
			elem := c.GetElementById(i.Id)
			if elem.Js.IsNull() || elem.Js.IsUndefined() {
				return fmt.Errorf("Element not found")
			}
			for event, handler := range i.Handlers {
				elem.Js.Call("addEventListener", event, handler)
			}
		} else {
			nodes := c.GetElementsByClassName(i.Class)
			for _, node := range nodes {
				for event, handler := range i.Handlers {
					node.Js.Call("addEventListener", event, handler)
				}
			}
		}
	}
	return nil
}

func (c *Client) Reload() error {
	for _, i := range c.Islands {
		if i.Id == "" && i.Class == "" {
			return fmt.Errorf("Element id or class not found")
		}

		if i.Id != "" {
			elem := c.GetElementById(i.Id)
			if elem.Js.IsNull() || elem.Js.IsUndefined() {
				return fmt.Errorf("Element not found")
			}
			for event, handler := range i.Handlers {
				elem.Js.Call("addEventListener", event, handler)
			}
		} else {
			nodes := c.GetElementsByClassName(i.Class)
			for _, node := range nodes {
				for event, handler := range i.Handlers {
					node.Js.Call("addEventListener", event, handler)
				}
			}
		}
	}
	return nil
}

func (c *Client) CreateElement(tag string) Node {
	return Node{Js: c.Dom.Call("createElement", tag)}
}

func (c *Client) GetElementById(id string) Node {
	return Node{Js: c.Dom.Call("getElementById", id)}
}

func (c *Client) GetElementsByClassName(class string) []Node {
	nodes := c.Dom.Call("getElementsByClassName", class)
	result := make([]Node, nodes.Length())
	for i := 0; i < nodes.Length(); i++ {
		result[i] = Node{Js: nodes.Index(i)}
	}
	return result
}

func (c *Client) QuerySelector(selector string) js.Value {
	return c.Dom.Call("querySelector", selector)
}

func (n *Node) SetInnerHTML(data string) {
	n.Js.Set("innerHTML", data)
}

func (n *Node) GetInnerHTML() string {
	return n.Js.Get("innerHTML").String()
}

func (n *Node) AppendChild(child Node) {
	n.Js.Call("appendChild", child.Js)
}
