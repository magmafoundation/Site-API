package utils

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v36/github"
	"golang.org/x/oauth2"
)

type VersionUtils struct {
	Client *github.Client
	RDB    *redis.Client
}

func (util *VersionUtils) Setup(token string) {
	c := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(c, ts)

	util.Client = github.NewClient(tc)

	util.RDB = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}

// GetStableReleases Fetch a list of all stable releases.
func (util *VersionUtils) GetStableReleases(repo string) []*github.RepositoryRelease {
	fmt.Printf("Fetching versions for %s\n", repo)

	val, err := util.RDB.Get(context.Background(), "releases:stable:"+repo).Result()

	if err == redis.Nil {
		fmt.Println("Fetching Stable releases.")

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
		if stableReleases == nil {
			stableReleases = []*github.RepositoryRelease{}
		}
		json, _ := json2.Marshal(releases)
		util.RDB.Set(context.Background(), "releases:stable:"+repo, json, 5*time.Minute)

		return stableReleases
	}

	var releases []*github.RepositoryRelease

	json2.Unmarshal([]byte(val), &releases)
	return releases
}

// GetPreReleases Fetch a list of all pre-releases.
func (util *VersionUtils) GetPreReleases(repo string) []*github.RepositoryRelease {

	val, err := util.RDB.Get(context.Background(), "releases:pre:"+repo).Result()

	if err == redis.Nil {
		fmt.Println("Fetching dev releases.")
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
		json, _ := json2.Marshal(preReleases)
		util.RDB.Set(context.Background(), "releases:pre:"+repo, json, 5*time.Minute)
		// Return the stable releases
		return preReleases
	}
	var preReleases []*github.RepositoryRelease

	json2.Unmarshal([]byte(val), &preReleases)

	if preReleases != nil {
		return preReleases
	}
	return []*github.RepositoryRelease{}

}
