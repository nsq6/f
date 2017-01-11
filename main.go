package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os/exec"
	"runtime"

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
	codeChannel, err = startWebServer(config.RedirectURL)
	if err != nil {
		fmt.Println(err)
	}
	accessType := oauth2.SetAuthURLParam("access_type", "offline")
	url := config.AuthCodeURL("dsadsadasdasd", accessType)
	openURL(url)
	code := <-codeChannel
	context := context.Background()
	token, err := config.Exchange(context, code)
	if err != nil {
		fmt.Println(err)
	}
	client := config.Client(context, token)
	resp, err := client.Get("https://www.googleapis.com/youtube/v3/channels")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp.Body)
}

func startWebServer(listenURI string) (codeCh chan string, err error) {
	codeCh = make(chan string)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := r.FormValue("code")
		codeCh <- code // send code to OAuth flow
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "Received code: %v\r\nYou can now safely close this browser window.", code)
	})
	go http.ListenAndServe(":9090", nil)
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
