package cmd

import (
	"fmt"
	"net"
	"net/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve dir",
	Short: "Serve a website from static files",
	Args:  cobra.ExactArgs(1),
	RunE:  runServe,
}

var port uint16

func init() {

	serveCmd.Flags().Uint16VarP(&port, "port", "p", port, "the port number to bind to")
	rootCmd.AddCommand(serveCmd)
}

func runServe(_ *cobra.Command, args []string) error {
	addr := fmt.Sprintf("localhost:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	http.Handle("/", http.FileServer(http.Dir(args[0])))

	fmt.Println("Server is running on http://" + listener.Addr().String())

	return http.Serve(listener, nil)
}
