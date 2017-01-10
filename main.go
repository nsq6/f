package main

import (
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

var (
	configFilePath = flag.String("config", "client_secret.live.json", "Path to client configuration file")
)

func main() {
	config, err := createConfig(configFilePath)
	if err != nil {
		fmt.Println(err)
	}
	var codeChannel chan string
	codeChannel, err = startHTTPWebServer(config.RedirectURL)
	if err != nil {
		fmt.Println(err)
	}
	accessType := oauth2.SetAuthURLParam("access_type", "offline")
	url := config.AuthCodeURL("dsadsadasdasd", accessType)
	openURL(url)
	code := <-codeChannel
	_, err = config.Exchange(context.Background(), code)
}

func startHTTPWebServer(listenURI string) (codeCh chan string, err error) {
	codeCh = make(chan string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
	})
	go func() {
		err = http.ListenAndServe(":9090", nil) // set listen port
		if err != nil {
			fmt.Printf("Cannot start web server on %v", listenURI)
			panic(err)
		}
	}()
	return codeCh, nil
}

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}
