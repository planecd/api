package clients

import (
	"log"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/mateothegreat/go-multilog/multilog"

	"github.com/google/go-github/v66/github"
)

type GitHubClient struct {
	Secret string
	Client *github.Client
}

type CallbackFn func(event *github.WorkflowRunEvent)

func (c *GitHubClient) Init(appID int64, installationID int64, keyPath string) error {
	// Shared transport to reuse TCP connections.
	tr := http.DefaultTransport

	// Wrap the shared transport for use with the app ID 1 authenticating with installation ID 99.
	itr, err := ghinstallation.NewKeyFromFile(tr, appID, installationID, keyPath)
	if err != nil {
		return err
	}

	// Use installation transport with github.com/google/go-github
	c.Client = github.NewClient(&http.Client{
		Transport: itr,
	})

	return nil
}

func (c *GitHubClient) Handle(req *http.Request, callback CallbackFn) error {
	payload, err := github.ValidatePayload(req, []byte(c.Secret))
	if err != nil {
		return err
	}

	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		return err
	}

	switch event := event.(type) {
	case *github.WorkflowJobEvent:
		multilog.Debug("webhook", "WorkflowJobEvent", map[string]any{
			"action": *event.Action,
		})
	case *github.WorkflowRunEvent:
		multilog.Debug("webhook", "WorkflowRunEvent", map[string]any{
			"action": *event.Action,
		})
	default:
		log.Println("Unsupported event type")
	}
	return nil
}
