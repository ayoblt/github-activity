package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Response struct {
	Type      string  `json:"type"`
	Repo      Repo    `json:"repo"`
	Payload   Payload `json:"payload"`
	CreatedAt string  `json:"created_at"`
}

type Repo struct {
	Name string `json:"name"`
}

type Payload struct {
	Action string `json:"action"`
	// Comment string `json:"comment"`
	Ref     string `json:"ref"`
	RefType string `json:"ref_type"`
	FullRef string `json:"full_ref"`
	Number  int    `json:"number"`
}

const timeFormat = "Mon Jan 2"

func main() {
	response, err := http.Get("http://api.github.com/users/ayoblt/events/public")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var responseObject []Response
	err = json.Unmarshal([]byte(body), &responseObject)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Output:")

	for _, item := range responseObject {
		t, _ := time.Parse(time.RFC3339, item.CreatedAt)
		switch item.Type {
		case "PushEvent":
			fmt.Println("- Pushed to", item.Repo.Name, "on", t.Format(timeFormat))
		case "CommitCommentEvent":
			fmt.Println("-", item.Payload.Action, "on a commit in", item.Repo.Name, "on", t.Format(timeFormat))
		case "CreateEvent":
			fmt.Println("- Created", item.Payload.RefType, item.Payload.Ref, "in", item.Repo.Name, "on", t.Format(timeFormat))
		case "DeleteEvent":
			fmt.Println("- Deleted", item.Payload.RefType, item.Payload.Ref, "in", item.Repo.Name, "on", t.Format(timeFormat))
		case "DiscussionEvent":
			fmt.Println("-", item.Payload.Action, "a discussion in", item.Repo.Name, "on", t.Format(timeFormat))
		case "ForkEvent":
			fmt.Println("- Forked", item.Repo.Name, "on", t.Format(timeFormat))
		case "GollumEvent":
			fmt.Println("- Edited wiki pages in", item.Repo.Name, "on", t.Format(timeFormat))
		case "IssueCommentEvent":
			fmt.Println("- Commented on an issue in", item.Repo.Name, "on", t.Format(timeFormat))
		case "PullRequestEvent":
			fmt.Println("-", item.Payload.Action, "a pull request in", item.Repo.Name, "on", t.Format(timeFormat))
		case "PullRequestReviewEvent":
			fmt.Println("- Reviewed a pull request in", item.Repo.Name, "on", t.Format(timeFormat))
		case "PullRequestReviewCommentEvent":
			fmt.Println("- Commented on a pull request in", item.Repo.Name, "on", t.Format(timeFormat))
		case "ReleaseEvent":
			fmt.Println("- Published a release in", item.Repo.Name, "on", t.Format(timeFormat))
		case "WatchEvent":
			fmt.Println("- Starred", item.Repo.Name, "on", t.Format(timeFormat))
		case "IssuesEvent":
			fmt.Println("-", item.Payload.Action, "an issue in", item.Repo.Name, "on", t.Format(timeFormat))
		default:
			fmt.Println("- Unregistered Event", item.Type)
		}
	}
}
