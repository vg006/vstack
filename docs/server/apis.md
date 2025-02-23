# Server APIs

The server package sets up and manages the web server, serving HTML templates, static files, and handling hot reloading.

## Functions
### NewRouter(dir string) *Router
  - Creates a server instance with a template directory.

  **Example:**
  ```go
  r := vstack.NewRouter("src")
  ```

## Methods
### (r *Router) Load() error
  - Loads all HTML templates from the directory into the server’s map.

  **Example:**
  ```go
  err := r.Load()
  if err != nil {
      fmt.Println("Error:", err)
  }
  ```

### (r *Router) InitHMR() error
  - Starts the server with HMR, watching for template changes.

  **Example:**
  ```go
  err := r.InitHMR()
  if err != nil {
      fmt.Println("Error:", err)
  }
  ```

### (r *Router) error
  - Reattaches event handlers after DOM changes (e.g., HMR updates).

  **Example:**
  ```go
  err := r.Reload()
  if err != nil {
      fmt.Println("Error:", err)
  }
  ```
