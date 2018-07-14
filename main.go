package main

import (
  "context"
  "flag"
  "log"
  "os"
  "regexp"

	"golang.org/x/oauth2"
  "github.com/google/go-github/github"
)

type app struct {
  Context   context.Context
  Client    *github.Client
  Delete    *bool
  Org       *string
  ExcludeRe *string
  IncludeRe *string
}

func newApp() *app {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &app{ Client: github.NewClient(tc), Context: ctx }
}

func (app *app) ListWatched() []*github.Repository {
  var allRepos []*github.Repository
  opts := &github.ListOptions{ PerPage: 100 }
  for {
    repos, resp, err := app.Client.Activity.ListWatched(app.Context, "", opts)
    if err != nil {
      log.Fatalf("unable to fully paginate subscriptions: %v", err)
    }
    allRepos = append(allRepos, repos...)
    if resp.NextPage == 0 {
      break
    }
    opts.Page = resp.NextPage
  }
  return allRepos
}

func (app *app) RepoInScope(repo *github.Repository) bool {
  include,ierr := regexp.MatchString(*app.IncludeRe,*repo.Name)
  if ierr != nil {
    log.Fatal("invalid inclusion regular expression")
  }
  exclude := false
  if *app.ExcludeRe != "" {
    matched,eerr := regexp.MatchString(*app.ExcludeRe,*repo.Name)
    if eerr != nil {
      log.Fatal("invalid exclusion regular expression")
    }
    exclude = matched
  }
  orgmatch := true
  if *app.Org != "" {
    orgmatch = *app.Org == *repo.Owner.Login
  }
  return include && ! exclude && orgmatch
}

func main() {
  app := newApp()
  app.Delete    = flag.Bool("delete", false, "delete matching repository subscriptions")
  app.Org       = flag.String("org", "", "constrain operations to repositories within this organization")
  app.ExcludeRe = flag.String("exclude", "", "regular expression with which to exclude subscriptions")
  app.IncludeRe = flag.String("include", ".", "regular expression with which to include subscriptions")
  flag.Parse()

	repos := app.ListWatched()
  for _,repo := range repos {
    if app.RepoInScope(repo) {
      if *app.Delete {
        _,err := app.Client.Activity.DeleteRepositorySubscription(
          app.Context, *repo.Owner.Login, *repo.Name,
        )
        fmt := "deleted repository subscription %v/%v\n"
        if err != nil {
          fmt = "unable to delete subscription for %v/%v"
        }
        log.Printf(fmt, *repo.Owner.Login, *repo.Name)
      } else {
        log.Printf("found repository subscription %v/%v\n", *repo.Owner.Login, *repo.Name)
      }
    }
  }
}
