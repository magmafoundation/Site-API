package utils

import (
	mgithub "MagmaAPI/github"
	mredis "MagmaAPI/redis"
	"context"
	json2 "encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v36/github"
)

type VersionUtils struct {
}

// GetStableReleases Fetch a list of all stable releases.
func (util *VersionUtils) GetStableReleases(repo string) []*github.RepositoryRelease {
	log.Printf("Fetching versions for %s\n", repo)

	val, err := mredis.RDB.Get(context.Background(), "releases:stable:"+repo).Result()

	if err == redis.Nil {
		log.Println("Fetching Stable releases.")

		// Fetch the data.
		releases, _, err := mgithub.Client.Repositories.ListReleases(context.TODO(), "magmafoundation", repo, &github.ListOptions{})
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
		mredis.RDB.Set(context.Background(), "releases:stable:"+repo, json, 5*time.Minute)

		return stableReleases
	}

	var releases []*github.RepositoryRelease

	json2.Unmarshal([]byte(val), &releases)
	return releases
}

// GetPreReleases Fetch a list of all pre-releases.
func (util *VersionUtils) GetPreReleases(repo string) []*github.RepositoryRelease {

	val, err := mredis.RDB.Get(context.Background(), "releases:pre:"+repo).Result()

	if err == redis.Nil {
		log.Println("Fetching dev releases.")
		// Fetch the data.
		releases, _, err := mgithub.Client.Repositories.ListReleases(context.TODO(), "magmafoundation", repo, &github.ListOptions{})
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
		mredis.RDB.Set(context.Background(), "releases:pre:"+repo, json, 5*time.Minute)
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
