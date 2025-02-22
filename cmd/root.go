package cmd

import (
	"os"

	"github.com/spf13/cobra"
	asset "github.com/vg006/vstack/internal/assets"
)

var rootCmd = &cobra.Command{
	Use:   "vstack",
	Short: "An innovative meta framework for seamless SSR and client interactivity",
	Long: `VStack is an innovative meta framework designed to seamlessly integrate 
server-side rendering with client-side interactivity.

It ensures real-time UI updates, centralized state, and dynamic data management 
capabilities by leveraging WebSockets and WebAssembly (Wasm).

With VStack, developers can efficiently develop and build client logic in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(asset.VStackLogo)
		cmd.Println(asset.Text.Render("Welcome to VStack!"))
		cmd.Println(asset.Text.Render("VStack is an innovative meta framework for seamless SSR and client interactivity.\n"))
		cmd.Println(asset.Text.Render("Use 'vstack --help' for more information."))
	},
}

// --- Executes the command ---
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
