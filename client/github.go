package client

import (
	"context"
	"os/exec"

	"github.com/google/go-github/v52/github"
	"golang.org/x/oauth2"
)

var pagination = 150

type IGithubClient interface {
	Repositories(org string) (*RepositoriesResult, error)
}

type githubclient struct {
	client *github.Client
}

type RepositoriesResult struct {
	Repositories []Repository
}

type Repository struct {
	Name     string
	Language string
	SSHUrl   string
}

func NewGithubClient(key string) IGithubClient {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: key},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &githubclient{
		client: client,
	}
}

func (ghc *githubclient) Repositories(org string) (*RepositoriesResult, error) {
	ctx := context.Background()

	opt := github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: pagination},
	}

	repos, _, err := ghc.client.Repositories.ListByOrg(ctx, org, &opt)

	if err != nil {
		return nil, err
	}

	reps := make([]Repository, len(repos))

	for i, repo := range repos {
		reps[i] = Repository{
			Name:     repo.GetName(),
			Language: repo.GetLanguage(),
			SSHUrl:   repo.GetSSHURL(),
		}
	}

	return &RepositoriesResult{
		Repositories: reps,
	}, nil
}

func (r *RepositoriesResult) FindRepoByName(name string) Repository {
	for _, v := range r.Repositories {
		if v.Name+" - "+v.Language == name {
			return v
		}
	}

	return Repository{}
}

func (r *RepositoriesResult) RepositoryNames() []string {
	names := make([]string, len(r.Repositories))

	for i, v := range r.Repositories {
		names[i] = v.Name + " - " + v.Language
	}

	return names
}

func (r *Repository) Clone() error {
	url := r.SSHUrl

	cmd := exec.Command("git", "clone", url)

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}
