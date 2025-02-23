# Vstack

## 📖 Description
> **TL;DR**: Vstack is a minimalistic frontend tooling that boosts the productivity of the Go developers in developing and building frontend for web applications using Go.

Vstack is a Go-based frontend tooling that enables developers to build dynamic web applications with minimal JavaScript.

By leveraging Wasm, Vstack provides a unified architecture for client-managed and server-driven rendering, with reactive state management and real-time updates.

## 🚀 [Demo](https://youtu.be/GnKUMbn62Nc1)

## 💡 Inspirations

- **[Htmx](https://htmx.org/):**
  For its minimalistic approach to dynamic web applications.

- **[Svelte](https://svelte.dev/):**
  For its reactive state management and efficient updates.

- **[Next.js](https://nextjs.org/):**
  For its file-based routing and hot module reloading.

This project is inspired by the above technologies and aims to provide a similar experience in an unified way for Go developers building frontend applications.

## ✨ Features
- **Unified API:**

  A single, high-level interface simplifies configuration and integration with existing Go projects.
- **Framework-Agnostic:**

  Developers can place templates in distinct directories and provide default data or middleware for each route.
- **Component Architecture:**

  Organize your templates into separate directories for client-managed and server-driven data, allowing precise control over which parts of your UI update dynamically.
- **Minimal JavaScript:**

  Use Go to manage your frontend state and updates, with minimal JavaScript required for client-side interactivity.


## 🎯 To-dos
- [ ] **Refine Architecture:**

  Define a clear separation between client-managed and server-driven rendering, with a unified API for both.
- [ ] **Client-Side Rendering:**

  Use WebAssembly to render templates and manage state on the client side.
- [ ] **Reactive State Management:**

  Define reactive variables ("runes") with minimal API exposure; updates trigger automatic UI refreshes.
- [ ] **Real-Time Updates:**

  Efficiently update only the parts of the UI that change, using a unified architecture that supports both client-managed and server-driven rendering (maybe using the SSE).

## ⬇️ Installation
  ```bash
  go get github.com/vg006/vstack@latest
  ```

## 🛠️ Example Usage
  To use Vstack in your project, follow the step:
  1. **Install the package:**:
      ```bash
      go get github.com/vg006/vstack@latest
      ```
  2. **Import the package**:
      ```go
      import "github.com/vg006/vstack"
      ```
  3. **Use the package**:
      1. **Client-side**:
          ```go
          package main

          import (
              "github.com/vg006/vstack/pkg"
          )

          func main() {
              c := vstack.NewClient()
              elem := c.GetByElementId("myDiv")
              elem.SetInnerHTML("Hello, World!")
          }
          ```
      2. **Server-side**:
          ```go
          package main

          import (
          	"github.com/vg006/vstack/pkg"
          )
          func main() {
          	s := vstack.NewServer("dir")
           	if err != nil {
            		fmt.Println("Error:", err)
          	}
           	http.ListenAndServe(":8080", s)
          }

  **Note**: The above code snippets are just examples. Kindly check the [docs](https://github.com/vg006/vstack/tree/main/docs) for more information on the APIs.

  5. **Build the WASM Client**:
      ```bash
      $ GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o dir/public/main.wasm path/to/client.go

      $ cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" dir/public/
      ```
  **Note**: Replace `dir` with the directory containing your client-side templates and `path/to/client.go` with the path to your client-side Go file.

  6. **Run the server**:
      ```bash
      go run cmd/server/main.go
      ```

## 🤝 Contributing
Contributions are welcome! Let's build, spread and contribute to Go community. Follow our coding standards and include tests for new features.

### Prerequisites
- Go 1.24.0 or later
- A modern web browser with WebAssembly support

### Steps
  1. Fork the repository.
  2. Create a new branch:
      ```bash
      git checkout -b feature/your-feature
      ```
  3. Commit your changes:
      ```bash
      git commit -m "Add your feature"
      ```
  4. Push the branch:
      ```bash
      git push origin feature/your-feature
      ```
  5. Open a pull request.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](https://github.com/vg006/vstack/tree/main?tab=MIT-1-ov-file) file for details.
