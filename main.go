package main

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v63/github"
)

func main() {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 959127, 53377449, "humop-checks-test.2024-08-01.private-key.pem")
	if err != nil {
		log.Fatal(err)
	}

	client := github.NewClient(&http.Client{
		Transport: itr,
	})

	ctx := context.Background()
	list, _, err := client.PullRequests.List(ctx, "ophum", "github-checks-test", &github.PullRequestListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	owner := "ophum"
	repo := "github-checks-test"
	for _, pr := range list {
		log.Println(*pr.Title, *pr.Base.Label)

		log.Println("base:", *pr.Base.Ref, *pr.Base.SHA)
		log.Println("head:", *pr.Head.Ref, *pr.Head.SHA)
		checkRunName := "test"
		checkRuns, _, err := client.Checks.ListCheckRunsForRef(ctx, owner, repo, pr.GetHead().GetSHA(), &github.ListCheckRunsOptions{})
		if err != nil {
			log.Fatal(err)
		}
		log.Println(checkRuns.GetTotal())
		index := slices.IndexFunc(checkRuns.CheckRuns, func(c *github.CheckRun) bool {
			log.Println(c.GetName())
			return c.GetName() == checkRunName
		})
		var checkRun *github.CheckRun
		if index == -1 {
			checkRun, _, err = client.Checks.CreateCheckRun(ctx, owner, repo, github.CreateCheckRunOptions{
				Name:    checkRunName,
				HeadSHA: *pr.Head.SHA,
			})
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Println("found")
			checkRun = checkRuns.CheckRuns[index]
		}

		log.Println("status:", checkRun.GetStatus())

		if _, _, err := client.Checks.UpdateCheckRun(ctx, owner, repo, checkRun.GetID(), github.UpdateCheckRunOptions{
			Name:       checkRunName,
			Conclusion: github.String("success"),
		}); err != nil {
			log.Fatal(err)
		}
	}
}
