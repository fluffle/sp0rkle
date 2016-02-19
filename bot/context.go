package bot

import (
	"fmt"
	"strings"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/sp0rkle/collections/conf"
	"github.com/fluffle/sp0rkle/util"
)

// Basic types representing the information we want to store about IRC things
type Nick string

func (n Nick) Lower() string {
	return strings.ToLower(string(n))
}

type Chan string

func (c Chan) Lower() string {
	return strings.ToLower(string(c))
}

// context encapsulates the bot's stuff and provides an interface
// for the rest of the bot to interact with IRC.
type Context struct {
	*client.Line
	Addressed bool

	conn *client.Conn
	rws  RewriteSet
}

func context(conn *client.Conn, line *client.Line) *Context {
	ctx := &Context{conn: conn, Line: line.Copy(), rws: bot.rewriters}
	// This is a bit of a dirty hack; context() returns nil to ignore a line.
	// TODO(fluffle): Ignores based on masks (or more likely regex).
	if conf.Ns(ignoreNs).String(strings.ToLower(ctx.Nick)) != "" {
		return nil
	}
	if ctx.Cmd != client.PRIVMSG {
		return ctx
	}
	ctx.Args[1], ctx.Addressed = util.RemovePrefixedNick(
		strings.TrimSpace(ctx.Args[1]), ctx.Me())
	// If we're being talked to in private, line.Args[0] will contain our Nick.
	// We should consider this as "addressing" us, and set Addressed = true
	if ctx.Args[0] == ctx.Me() {
		ctx.Addressed = true
	}
	return ctx
}

func (ctx *Context) Storable() (Nick, Chan) {
	return Nick(ctx.Nick), Chan(ctx.Args[0])
}

// ReplyN() adds a prefix of "nick: " to the reply text,
func (ctx *Context) ReplyN(fm string, args ...interface{}) {
	args = append([]interface{}{ctx.Nick}, args...)
	ctx.Reply("%s: "+fm, args...)
}

// whereas Reply() does not.
func (ctx *Context) Reply(fm string, args ...interface{}) {
	ctx.conn.Privmsg(ctx.Target(),
		ctx.rws.Rewrite(fmt.Sprintf(fm, args...), ctx))
}

func (ctx *Context) Do(fm string, args ...interface{}) {
	ctx.conn.Action(ctx.Target(),
		ctx.rws.Rewrite(fmt.Sprintf(fm, args...), ctx))
}

func (ctx *Context) Privmsg(ch, text string) {
	ctx.conn.Privmsg(ch, text)
}

func (ctx *Context) Action(ch, text string) {
	ctx.conn.Action(ch, text)
}

func (ctx *Context) Topic(ch string, topic ...string) {
	ctx.conn.Topic(ch, topic...)
}

func (ctx *Context) Me() string {
	return ctx.conn.Me().Nick
}
