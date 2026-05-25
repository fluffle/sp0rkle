package bot

import (
	"flag"
	"strings"

	"github.com/fluffle/golog/logging"
)

var (
	channels *string = flag.String("channels", "#sp0rklf",
		"Comma-separated list of channels to join.")
	rebuilder *string = flag.String("rebuilder", "",
		"Nick[:password] to accept rebuild command from.")
	oper *string = flag.String("oper", "",
		"user:password for server OPER command on connect, or $ENV_VAR or <file_path to secret.")
	vhost *string = flag.String("vhost", "",
		"user:password for server VHOST command on connect, or $ENV_VAR or <file_path to secret.")
)

func (b *Bot) connected(ctx *Context) {
	ctx.conn.Mode(ctx.Me(), "+B")
	if GetSecret(*oper) != "" {
		up := strings.SplitN(*oper, ":", 2)
		if len(up) == 2 {
			ctx.conn.Oper(up[0], up[1])
		}
	}
	if GetSecret(*vhost) != "" {
		up := strings.SplitN(*vhost, ":", 2)
		if len(up) == 2 {
			ctx.conn.VHost(up[0], up[1])
		}
	}
	for _, c := range strings.Split(*channels, ",") {
		logging.Info("Joining %s on startup.\n", c)
		ctx.conn.Join(c)
	}
}

func (b *Bot) shutdown(ctx *Context) {
	s := strings.Split(GetSecret(*rebuilder), ":")
	if s[0] == "" || s[0] != ctx.Nick || !strings.HasPrefix(strings.ToLower(ctx.Text()), "shutdown") {
		return
	}
	fields := strings.Fields(ctx.Text())
	if len(s) > 1 && fields[len(fields)-1] != s[1] {
		return
	}
	b.servers.Shutdown(false)
}
