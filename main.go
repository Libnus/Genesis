package main

import (
    "log"
    "os/exec"
    "runtime"
)

func main() {
    // Directly open the browser to the existing server at port 8000
    openBrowser("http://localhost:8000/brain.html")
}

// openBrowser opens the default browser to a specific URL
func openBrowser(url string) {
    var err error
    switch runtime.GOOS {
    case "linux":
        err = exec.Command("xdg-open", url).Start()
    case "windows":
        err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
    case "darwin":
        err = exec.Command("open", url).Start()
    }
    if err != nil {
        log.Printf("Failed to open browser: %v", err)
    }
}
