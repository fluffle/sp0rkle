package netdriver

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/util"
	"github.com/google/go-github/github"
	"strings"
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

func sp(s string) *string {
	//  FFFUUUUuuu string pointers in Issue literals.
	return &s
}

func githubClient() *github.Client {
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: *githubToken}}
	return github.NewClient(t.Client())
}

func githubCreateIssue(ctx *bot.Context) {
	s := strings.SplitN(ctx.Text(), ". ", 2)
	if s[0] == "" {
		ctx.ReplyN("I'm not going to create an empty issue.")
		return
	}

	issue := &github.Issue{Title: sp(s[0] + ".")}
	if len(s) == 2 {
		issue.Body = &s[1]
	}
	issue, _, err := gh.Issues.Create(githubUser, githubRepo, issue)
	if err != nil {
		ctx.ReplyN("Error creating issue: %v", err)
		return
	}

	// Can't set labels on create due to go-github #75 :/
	_, _, err = gh.Issues.ReplaceLabelsForIssue(
		githubUser, githubRepo, *issue.Number,
		[]string{"from:IRC", "nick:" + ctx.Nick, "chan:" + ctx.Target()})
	if err != nil {
		ctx.ReplyN("Failed to add labels to issue: %v", err)
	}
	ctx.ReplyN("Issue #%d created at %s/%d",
		*issue.Number, githubIssuesURL, *issue.Number)
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

	labels, _, err := gh.Issues.ListLabelsByIssue(githubUser, githubRepo, issue)
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
