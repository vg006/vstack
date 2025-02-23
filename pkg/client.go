//go:build js && wasm

package vstack

import (
	"syscall/js"

	"github.com/vg006/vstack/internal/client"
)

type (
	Client = client.Client
	Node   = client.Node
	Island = client.Island
)


func NewClient() Client {
	return Client{Dom: js.Global().Get("document")}
}
