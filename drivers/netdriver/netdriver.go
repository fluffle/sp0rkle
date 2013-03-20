package netdriver

import (
	"flag"
	"github.com/fluffle/sp0rkle/bot"
	"io/ioutil"
	"net/http"
)

var githubToken = flag.String("github_token", "",
	"OAuth2 token for accessing the GitHub API.")

func get(req string) ([]byte, error) {
	res, err := http.Get(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func Init() {
	bot.Command(urbanDictionary, "ud", "ud <term>  -- "+
		"Look up <term> on UrbanDictionary.")
	bot.Command(createGitHubIssue, "file bug:", "file bug: <title>. "+
		"<descriptive body>  -- Files a bug on GitHub. Abusers will be hurt.")
}
