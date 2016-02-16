package netdriver

import (
	"flag"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util"
	"github.com/google/go-github/github"
)

var (
	githubToken = flag.String("github_token", "",
		"OAuth2 token for accessing the GitHub API.")
)

const (
	githubUser      = "fluffle"
	githubRepo      = "sp0rkle"
	githubURL       = "https://github.com/" + githubUser + "/" + githubRepo
	githubIssuesURL = githubURL + "/issues"
)

func githubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: bot.GetSecret(*githubToken)},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}

func githubCreateIssue(ctx *bot.Context) {
	s := strings.SplitN(ctx.Text(), ". ", 2)
	if s[0] == "" {
		ctx.ReplyN("I'm not going to create an empty issue.")
		return
	}

	req := &github.IssueRequest{
		Title: &s[0],
		Labels: &[]string{
			"from:IRC",
			"nick:" + ctx.Nick,
			"chan:" + ctx.Target(),
		},
	}
	if len(s) == 2 {
		req.Body = &s[1]
	}
	issue, _, err := gh.Issues.Create(githubUser, githubRepo, req)
	if err != nil {
		ctx.ReplyN("Error creating issue: %v", err)
		return
	}

	ctx.ReplyN("Issue #%d created at %s/%d",
		*issue.Number, githubIssuesURL, *issue.Number)
}

func githubUpdateIssue(ctx *bot.Context) {
	l := &util.Lexer{Input: ctx.Text()}
	issue := int(l.Number())
	if issue == 0 {
		ctx.ReplyN("Not sure what issue you're talking about?")
		return
	}
	text := strings.TrimSpace(l.Find(0))
	if text == "" {
		ctx.ReplyN("Don't you have anything to say?")
		return
	}
	text = fmt.Sprintf("<%s/%s> %s", ctx.Nick, ctx.Target(), text)
	comm, _, err := gh.Issues.CreateComment(
		githubUser, githubRepo, issue, &github.IssueComment{Body: &text})
	if err != nil {
		ctx.ReplyN("Error creating issue comment: %v", err)
		return
	}
	ctx.ReplyN("Created comment %s/%d#issuecomment-%d",
		githubIssuesURL, issue, *comm.ID)
}

func githubWatcher(ctx *bot.Context) {
	// Watch #sp0rklf for IRC messages about issues coming from github.
	if ctx.Nick != "fluffle\\sp0rkle" || ctx.Target() != "#sp0rklf" ||
		!strings.Contains(ctx.Text(), "issue #") {
		return
	}

	text := util.RemoveColours(ctx.Text()) // srsly github why colours :(
	l := &util.Lexer{Input: text}
	l.Find(' ')
	text = text[l.Pos()+1:]
	l.Find('#')
	l.Next()
	issue := int(l.Number())

	labels, _, err := gh.Issues.ListLabelsByIssue(
		githubUser, githubRepo, issue, &github.ListOptions{})
	if err != nil {
		logging.Error("Error getting labels for issue %d: %v", issue, err)
		return
	}
	for _, l := range labels {
		kv := strings.Split(*l.Name, ":")
		if len(kv) == 2 && kv[0] == "nick" {
			logging.Debug("Recording tell for %s about issue %d.", kv[1], issue)
			r := reminders.NewTell("that "+text, bot.Nick(kv[1]), "github", "")
			if err := rc.Insert(r); err != nil {
				logging.Error("Error inserting github tell: %v", err)
			}
		}
	}
}
