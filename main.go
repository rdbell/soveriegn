package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/antonholmquist/jason"
	"github.com/kr/pretty"
)

func main() {
	// Use keychain path/password as stored in the environment vars
	keychainPath := os.Getenv("KEYCHAIN_PATH")
	keychainPassword := os.Getenv("KEYCHAIN_PASSWORD")

	// Read search query from command line argument
	query := os.Args[1]

	// Extract the keychain using the 1password.js logic
	// TODO: Rewrite 1password.js in Go
	cmd := exec.Command("/bin/sh", "-c", "node ./1password.js \""+keychainPath+"\" \""+keychainPassword+"\"")
	out, err := cmd.Output()

	// Read the keychain items using the Jason library
	rawKeychain := string(out)
	rawKeychain = rawKeychain[0 : len(rawKeychain)-2]
	formattedKeychain := "{\"keychain\" : [" + rawKeychain + "]}"

	if err != nil {
		fmt.Println("Error block 1")
		panic(err)
	}

	keychainItems, err := jason.NewObjectFromBytes([]byte(formattedKeychain))

	if err != nil {
		fmt.Println("Error block 2")
		panic(err)
	}

	keychain, err := keychainItems.GetObjectArray("keychain")

	if err != nil {
		fmt.Println("Error block 3")
		panic(err)
	}

	var results []*jason.Object

	// Loop through the keychain to find a URL matching the user's query
	for _, keychainItem := range keychain {
		urls, err := keychainItem.GetObjectArray("URLs")
		if err == nil {
			for _, urlEntry := range urls {
				url, err := urlEntry.GetString("url")
				if err == nil {
					if strings.Contains(url, query) {
						results = append(results, keychainItem)
					}
				}
			}
		}
	}

	pretty.Println(results)
}
