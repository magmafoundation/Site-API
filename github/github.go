package mgithub

import "C"
import (
	"context"
	GH "github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

var (
	Client *GH.Client
)

func Setup(token string)  {
	c := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(c, ts)

	Client = GH.NewClient(tc)
}