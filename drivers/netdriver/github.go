package netdriver

import (
	"flag"
	"github.com/fluffle/go-github-client/client"
	"github.com/fluffle/go-github-client/issues"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"strings"
)

var githubToken = flag.String("github_token", "",
	"OAuth2 token for accessing the GitHub API.")

func createGitHubIssue(ctx *bot.Context) {
	if *githubToken == "" {
		ctx.ReplyN("I don't have a GitHub API token, sorry.")
		return
	}
	s := strings.SplitN(ctx.Text(), ". ", 2)
	if s[0] == "" {
		ctx.ReplyN("I'm not going to create an empty issue.")
		return
	}
	ghc, _ := client.NewGithubClient("", *githubToken, client.AUTH_OAUTH2_TOKEN)
	ic := issues.NewIssues(ghc)
	issue := &issues.IssueDataCreate{
		Title:    s[0]+".",
		Assignee: "fluffle",
		Labels:   []string{"from:IRC", "nick:"+ctx.Nick, "chan:"+ctx.Target()},
	}
	if len(s) == 2 {
		issue.Body = s[1]
	}
	res, err := ic.CreateIssue("fluffle", "sp0rkle", issue)
	if err != nil {
		ctx.ReplyN("Error creating issue: %v", err)
		return
	}
	data, err := res.JsonMap()
	if err != nil {
		ctx.ReplyN("Error unmarshalling JSON response: %v", err)
		return
	}
	logging.Info("%#v", data)
	ctx.ReplyN("Issue #%d created at %s",
		data.GetInt("number"), data.GetString("html_url"))
}
