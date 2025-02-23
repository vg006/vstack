## Client APIs

The client package provides APIs for DOM manipulation and event handling in the browser using WebAssembly (WASM). Below is a detailed list of methods and functions, with examples for ease of use.

## Functions
### NewClient() Client
  - Creates a new client instance.

  **Example:**
  ```go
  c := vstack.NewClient()
  ```

## Methods
### (c *Client) AddIslands(is ...Island) error
  - Adds islands, attaching event handlers to DOM elements.

  **Example:**
  ```go
  err := c.AddIslands(server.Island{Id: "delete", Handlers: map[string]js.Function{...}})
  if err != nil {
      fmt.Println("Error:", err)
  }
  ```

### (c *Client) GetElementById(id string) Node
  - Retrieves a DOM element by its ID.

  **Example:**
  ```go
  elem := c.GetElementById("myDiv")
  ```

### (c *Client) GetElementsByClassName(class string) []Node
  - Retrieves all elements with a specified class name.

  **Example:**
  ```go
  nodes := c.GetElementsByClassName("delete-me")
  for _, node := range nodes {
      // handle node
  }
  ```

### (c *Client) CreateElement(tagName string) Node
  - Creates a new DOM element with the given tag name.

  **Example:**
  ```go
  newDiv := c.CreateElement("div")
  newDiv.SetInnerHTML("Hello, World!")
  ```

### (n *Node) AppendChild(child Node)
  - Appends a child node to a parent node.

  **Example:**
  ```go
  todoList := c.GetElementById("todoList")
  newItem := c.CreateElement("div")
  todoList.AppendChild(newItem)
  ```

### (n *Node) RemoveChild(child Node)
  - Removes a child node from its parent.

  **Example:**
  ```go
  parent := c.GetByElementId("todoList")
  child := c.GetByElementId("item1")
  parent.RemoveChild(child)
  ```

### (n *Node) AddEventListener(event string, handler js.Func)
  - Adds an event listener to a DOM element.

  **Example:**
  ```go
  button := c.GetByElementId("addTodo")
  button.AddEventListener("click", func(this js.Value, args []js.Value) interface{} {
      // handle click logic
      return nil
  })
  ```
### (n *Node) SetInnerHTML(html string)
  - Sets the inner HTML of a DOM element.

  **Example:**
  ```go
  div := c.GetElementById("myDiv")
  div.SetInnerHTML("<p>Hello, World!</p>")
  ```

### (n *Node) GetInnerHTML() string
  - Retrieves the inner HTML of a DOM element.

  **Example:**
  ```go
  div := c.GetElementById("myDiv")
  htmlContent := div.GetInnerHTML()
  fmt.Println(htmlContent)
  ```
