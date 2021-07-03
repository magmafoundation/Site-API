package utils

import (
	"context"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

type VersionUtils struct {
	Client *github.Client
}

func (util *VersionUtils) Setup(token string) {
	c := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(c, ts)

	util.Client = github.NewClient(tc)
}

// GetStableReleases Fetch a list of all stable releases.
func (util *VersionUtils) GetStableReleases(repo string) []*github.RepositoryRelease {

	// Fetch the data.
	releases, _, err := util.Client.Repositories.ListReleases(context.TODO(), "magmafoundation", repo, &github.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var stableReleases []*github.RepositoryRelease

	// Filter out all pre-releases
	for _, release := range releases {
		if !*release.Prerelease {
			stableReleases = append(stableReleases, release)
		}
	}

	// Return the stable releases
	return stableReleases

}

// GetPreReleases Fetch a list of all pre-releases.
func (util *VersionUtils) GetPreReleases(repo string) []*github.RepositoryRelease {
	// Fetch the data.
	releases, _, err := util.Client.Repositories.ListReleases(context.TODO(), "magmafoundation", repo, &github.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	var preReleases []*github.RepositoryRelease

	// Filter out all pre-releases
	for _, release := range releases {
		if *release.Prerelease {
			preReleases = append(preReleases, release)
		}
	}

	// Return the stable releases
	return preReleases
}


