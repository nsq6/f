package main

import (
	"golang.org/x/oauth2"
)

type TokenSource struct {
	oauth2.TokenSource
}
