package main

import (
	"context"

	"encoding/json"
	"fmt"
	"log"
	"os"


	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	dropboxAppID := os.Getenv("DROPBOX_APP_ID")
	dropboxSecret := os.Getenv("DROPBOX_APP_SECRET")
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     dropboxAppID,
		ClientSecret: dropboxSecret,
		Scopes:       []string{"files.metadata.read", "files.content.read"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.dropbox.com/oauth2/authorize",
			TokenURL: "https://api.dropboxapi.com/oauth2/token",
		},
	}

	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	url := conf.AuthCodeURL("state", oauth2.SetAuthURLParam("token_access_type", "offline"), oauth2.AccessTypeOffline)

	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Printf("Once you have a code, past it and apress enter: ")

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(&tok)
	if err != nil {
		log.Fatal(err)
	}

	// client := conf.Client(ctx, tok)
	// resp, err := client.Post("https://api.dropboxapi.com/2/files/list_folder", "application/json", strings.NewReader(`{
	// 	"path": ""
	// }`))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer resp.Body.Close()
	// io.Copy(os.Stdout, resp.Body)
	//
	// err = json.NewEncoder(os.Stdout).Encode(&tok)
	// if err != nil {
	// 	log.Fatal(err)
	// }

}
