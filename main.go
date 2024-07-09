package main

import (
	"embed"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"

	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/routes"
)

// Embed the templates and static files
//
//go:embed templates/*
//go:embed static/*
var TemplateFiles embed.FS

func main() {
	file.VerificaPasta()
	// Check if the server is already running
	if isServerRunning("localhost:8000") {
		// Open the browser if the server is already running
		openBrowser("http://localhost:8000")
	} else {
		// Start the server and open the browser
		go func() {
			openBrowser("http://localhost:8000")
		}()
		routes.CarregaRotas(TemplateFiles)
		http.ListenAndServe(":8000", nil)
	}
}

// isServerRunning checks if a TCP server is running on the given address
func isServerRunning(address string) bool {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// openBrowser opens the default web browser at the given URL
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Printf("failed to open browser: %v\n", err)
	}
}
