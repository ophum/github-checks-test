package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/google/go-github/v63/github"
)

func main() {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, 1959127, 53377449, "humop-checks-test.2024-08-01.private-key.pem")
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

	for _, pr := range list {
		log.Println(pr.Title, pr.Base)
		j, _ := json.MarshalIndent(pr, "", "  ")
		log.Println(string(j))
	}
}
