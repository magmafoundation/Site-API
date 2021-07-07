package utils

import (
	mgithub "MagmaAPI/github"
	"context"
	"github.com/google/go-github/v36/github"
)

type Stats struct {
	ClosedIssues int `json:"closed_issues"`
	OpenIssues   int `json:"open_issues"`
	PlayerCount  int `json:"player_count"`
	ServerCount  int `json:"server_count"`
}

func GetClosedIssues(repo string) int {

	opt := &github.IssueListByRepoOptions{
		State:       "closed",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allIssues []*github.Issue

	for {
		issues, resp, err := mgithub.Client.Issues.ListByRepo(context.Background(), "magmafoundation", repo, opt)
		if err != nil {
			panic(err.Error())
		}

		for _, issue := range issues {
			if !issue.IsPullRequest() {
				allIssues = append(allIssues, issue)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return len(allIssues)
}
func GetOpenIssues(repo string) int {

	opt := &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	var allIssues []*github.Issue

	for {
		issues, resp, err := mgithub.Client.Issues.ListByRepo(context.Background(), "magmafoundation", repo, opt)
		if err != nil {
			panic(err.Error())
		}

		for _, issue := range issues {
			if !issue.IsPullRequest() {
				allIssues = append(allIssues, issue)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return len(allIssues)
}

func GetStats() Stats {

	stats := Stats{
		ClosedIssues: 0,
		OpenIssues:   0,
		PlayerCount:  0,
		ServerCount:  0,
	}

	var serversArray [][]int
	GetJson("https://bstats.org/api/v1/plugins/5445/charts/servers/data?maxElements=1", &serversArray)
	stats.ServerCount = serversArray[0][1]

	var playersArray [][]int
	GetJson("https://bstats.org/api/v1/plugins/5445/charts/players/data?maxElements=1", &playersArray)
	stats.PlayerCount = playersArray[0][1]

	repos := []string{"magma", "Magma-1.16.x"}

	for _, r := range repos {
		stats.ClosedIssues = stats.ClosedIssues + GetClosedIssues(r)
		stats.OpenIssues = stats.OpenIssues + GetOpenIssues(r)
	}

	return stats

}
