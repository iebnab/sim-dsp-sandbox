package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
)

func serveOK(response http.ResponseWriter, _ *http.Request) {
	response.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = response.Write([]byte("ok\n"))
}

func run(args []string, port string) error {
	if len(args) == 1 && args[0] == "--job" {
		return nil
	}
	if len(args) != 0 {
		return fmt.Errorf("unsupported fixture mode")
	}
	if port == "" {
		port = "8080"
	}
	listener, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	return http.Serve(listener, http.HandlerFunc(serveOK))
}

func main() {
	if err := run(os.Args[1:], os.Getenv("PORT")); err != nil {
		fmt.Fprintln(os.Stderr, "fixture startup failed:", err)
		os.Exit(1)
	}
}
