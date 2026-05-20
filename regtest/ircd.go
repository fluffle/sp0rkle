package regtest

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"text/template"
)

type TemplateContext struct {
	Instance string
	Address string
	TempDir string
}

func freeLocalAddr() (string, error) {
	// tyvm https://playtechnique.io/blog/finding-a-free-port-in-go.html
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return "", fmt.Errorf("could not open listen socket: %w", err)
	}
	defer listener.Close()
	return listener.Addr().String(), nil
}

func templateConfig(tmpDir, localAddr string) (string, error) {
	tc := TemplateContext{
		Instance: randSuffix(),
		Address: localAddr,
		TempDir: tmpDir,
	}
	configFile := fmt.Sprintf("ircd-%s.yaml", tc.Instance)
	abs, err := filepath.Abs(filepath.Join(tmpDir, configFile))
	if err != nil {
		return "", fmt.Errorf("determine abspath: %w", err)
	}
	b := bytes.NewBuffer(make([]byte, 0, len(ircdYAML)+128))
	if err = configTemplate.Execute(b, tc); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	if err = os.WriteFile(abs, b.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("write config: %w", err)
	}
	return abs, nil
}

var configTemplate = template.Must(
	template.New("ircd.yaml").Parse(ircdYAML))

const ircdYAML = `
network:
    name: "{{ .Instance }}"

server:
    name: "{{ .Instance }}.test"
    max-sendq: 96k
    listeners:
        "{{ .Address }}":
    lookup-hostnames: false
    forward-confirm-hostnames: false
    check-ident: false
    motd: "/dev/null"

    idle-timeouts:
        registration: 60s
        ping: 1m30s
        disconnect: 2m30s

    relaymsg:
        enabled: false

    ip-limits:
        count: false
        max-concurrent-connections: 64

        throttle: false
        window: 10m
        max-connections-per-window: 64

        cidr-len-ipv4: 32
        cidr-len-ipv6: 64

        exempted:
            - "localhost"

    ip-cloaking:
        enabled: false
    #output-path: "/home/ergo/out"


accounts:
    default-user-modes: +i
    authentication-enabled: false
    skip-server-password: false

    registration:
        enabled: false

    login-throttling:
        enabled: false

    nick-reservation:
        enabled: false

    multiclient:
        enabled: false

    vhosts:
        enabled: false

    auth-script:
        enabled: false

    oauth2:
        enabled: false

    jwt-auth:
        enabled: false

channels:
    default-modes: +ntC
    max-channels-per-client: 100
    operator-only-creation: false

    registration:
        enabled: false

logging:
    -
        method: stderr
        # listeners @ info level is necessary for the harness to
        # detect server startup by watching STDERR
        type: "listeners internal"

        # one of: debug info warn error
        level: info

debug:
    recover-from-errors: true

lock-file: "{{ .TempDir }}/{{ .Instance }}.lock"

datastore:
    path: "{{ .TempDir }}/{{ .Instance }}.db"
    autoupgrade: true

    mysql:
        enabled: false

    postgresql:
        enabled: false

    sqlite:
        enabled: false

languages:
    enabled: false

limits:
    nicklen: 32
    identlen: 20
    realnamelen: 150
    channellen: 64
    awaylen: 390
    kicklen: 390
    topiclen: 390
    monitor-entries: 100
    whowas-entries: 100
    chan-list-modes: 100
    registration-messages: 1024
    multiline:
        max-bytes: 4096 # 0 means disabled
        max-lines: 100  # 0 means no limit

fakelag:
    enabled: false

roleplay:
    enabled: false

history:
    enabled: false
    persistent:
        enabled: false

allow-environment-overrides: true

metadata:
    enabled: false

webpush:
    enabled: false

api:
    enabled: false
`

