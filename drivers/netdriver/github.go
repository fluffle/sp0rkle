package netdriver

import (
	"code.google.com/p/goauth2/oauth"
	"flag"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"strings"
	"time"
)

var (
	githubToken = flag.String("github_token", "",
		"OAuth2 token for accessing the GitHub API.")
	githubPollFreq = flag.Duration("github_poll_freq", 4 * time.Hour,
		"Frequency to poll github for bug updates.")
)

const (
	githubIssuesURL = "https://github.com/fluffle/sp0rkle/issues"
	ISO8601 = "2006-01-02T15:04:05Z"
)

func sp(s string) *string {
	//  FFFUUUUuuu string pointers in Issue literals.
	return &s
}

func githubClient() *github.Client {
	t := &oauth.Transport{Token: &oauth.Token{AccessToken: *githubToken}}
	return github.NewClient(t.Client())
}

func githubCreateIssue(ctx *bot.Context, gh *github.Client) {
	s := strings.SplitN(ctx.Text(), ". ", 2)
	if s[0] == "" {
		ctx.ReplyN("I'm not going to create an empty issue.")
		return
	}

	issue := &github.Issue{
		Title:    sp(s[0] + "."),
		Assignee: &github.User{Login: sp("fluffle")},
		Labels:   []github.Label{
			{Name: sp("from:IRC")},
			{Name: sp("nick:"+ctx.Nick)},
			{Name: sp("chan:"+ctx.Target())},
		},
	}
	if len(s) == 2 {
		issue.Body = &s[1]
	}
	issue, _, err := gh.Issues.Create("fluffle", "sp0rkle", issue)
	if err != nil {
		ctx.ReplyN("Error creating issue: %v", err)
		return
	}
	ctx.ReplyN("Issue #%d created at %s/%d",
		*issue.Number, githubIssuesURL, *issue.Number)
}

type ghPoller struct {
	// essentially a github client.
	*github.Client
}

type ghUpdate struct {
	issue                int
	updated, closed      time.Time
	nick, channel, title string
	comment, commenter   string
}

func (u ghUpdate) String() string {
	s := []string{}
	if !u.closed.IsZero() {
		s = append(s, fmt.Sprintf("Issue %s/%d (%s) closed at %s.",
			githubIssuesURL, u.issue, u.title, u.closed))
	} else {
		s = append(s, fmt.Sprintf("Issue %s/%d (%s) updated at %s.",
			githubIssuesURL, u.issue, u.title, u.updated))
	}
	if u.comment != "" {
		comment := u.comment
		trunc := " "
		if len(comment) > 100 {
			idx := strings.Index(comment, " ") + 100
			if idx >= 100 {
				comment = comment[:idx] + "..."
				trunc = " (truncated) "
			}
		}
		s = append(s, fmt.Sprintf("Recent%scomment by %s: %s",
			trunc, u.commenter, comment))
	}
	return strings.Join(s, " ")
}

func githubPoller(gh *github.Client) *ghPoller {
	return &ghPoller{gh}
}

func (ghp *ghPoller) Poll([]*bot.Context) { ghp.getIssues() }
func (ghp *ghPoller) Start() { /* empty */ }
func (ghp *ghPoller) Stop() { /* empty */ }
func (ghp *ghPoller) Tick() time.Duration { return *githubPollFreq }

func (ghp *ghPoller) getIssues() {
	opts := &github.IssueListByRepoOptions{
		Labels: []string{"from:IRC"},
		Sort:   "updated",
		State:  "open",
		Since:  time.Now().Add(*githubPollFreq * -1),
	}
	open, _, err := ghp.Issues.ListByRepo("fluffle", "sp0rkle", opts)
	if err != nil {
		logging.Error("Error listing open issues: %v", err)
	}
	opts.State = "closed"
	closed, _, err := ghp.Issues.ListByRepo("fluffle", "sp0rkle", opts)
	if err != nil {
		logging.Error("Error listing open issues: %v", err)
	}
	open = append(open, closed...)
	if len(open) == 0 { return }
	for _, issue := range open {
		update := ghp.parseIssue(issue)
		logging.Info("Adding tell for %s regarding issue %d.", update.nick, update.issue)
		r := reminders.NewTell(update.String(), bot.Nick(update.nick),
			"github", bot.Chan(update.channel))
		if err := rc.Insert(r); err != nil {
			logging.Error("Error inserting github tell: %v", err)
		}
	}
}

func (ghp *ghPoller) parseIssue(issue github.Issue) ghUpdate {
	update := ghUpdate{
		issue:   *issue.Number,
		updated: *issue.UpdatedAt,
		title:   *issue.Title,
	}
	if issue.ClosedAt != nil && time.Now().Sub(*issue.ClosedAt) < *githubPollFreq {
		update.closed = *issue.ClosedAt
	}
	for _, l := range issue.Labels {
		kv := strings.Split(*l.Name, ":")
		if len(kv) != 2 {
			continue
		}
		switch kv[0] {
		case "nick":
			update.nick = kv[1]
		case "chan":
			update.channel = kv[1]
		}
	}
	if *issue.Comments == 0 { return update }
	opts := &github.IssueListCommentsOptions{
		Sort: "updated",
		Direction: "desc",
		Since:  time.Now().Add(*githubPollFreq * -1),
	}
	comm, _ , err := ghp.Issues.ListComments(
		"fluffle", "sp0rkle", *issue.Number, opts)
	if err != nil {
		logging.Error("Error getting comments for issue %d: %v",
			*issue.Number, err)
	} else {
		update.comment = *comm[0].Body
		update.commenter = *comm[0].User.Login
	}
	return update
}
