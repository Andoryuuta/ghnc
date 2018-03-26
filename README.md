# ghnc (GitHub Name Checker)
[![GoDoc](https://godoc.org/github.com/Andoryuuta/ghnc?status.svg)](https://godoc.org/github.com/Andoryuuta/ghnc)

A GitHub username checking library for Go.

## Installation
`go get github.com/Andoryuuta/ghnc`

## Usage
```Go
package main

import (
	"fmt"

	"github.com/Andoryuuta/ghnc"
)

func main() {
	username := "Andoryuuta"

	// Get the client
	client, err := ghnc.GetGHClient()
	if err != nil {
		panic(err)
	}

	// Check if the username is available
	available, reason, err := client.UsernameAvailable(username)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Username %s available: %v (reason: %v)", username, available, reason)
}
```



