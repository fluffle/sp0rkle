[![Build Status](https://api.travis-ci.org/fluffle/sp0rkle.svg)](https://travis-ci.org/fluffle/sp0rkle)

Getting started, from scratch:

1) Install some dependencies and MANY version control systems.

 ```bash
sudo apt-get install build-essential bison mongodb libsqlite3-dev
sudo apt-get install git bzr

Ensure you have mongodb 2.x or higher
mongod --version
db version v2.2.1, pdfile version 4.5
```

If not the following page may help
https://docs.mongodb.org/manual/administration/install-on-linux/

2) Build go.

Install go https://golang.org/doc/install

 ```bash
# consider putting these in ~/.bashrc too...
export GOROOT="$HOME/go"
export GOPATH="$HOME/gocode"
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

# ... and creating the GOPATH directory structure.
# (read `go help gopath` for details of this)
mkdir -p $GOPATH/{src,pkg,bin}
```
3) Use the `go` tool to get dependencies.
 ```bash
go get github.com/fluffle/goirc/client
go get github.com/fluffle/golog/logging
go get github.com/kuroneko/gosqlite3
go get github.com/google/go-github/github
go get golang.org/x/oauth2
go get gopkg.in/mgo.v2
```
4) Clone sp0rkle's code from github.

 ```bash
cd $GOPATH/src/github.com/fluffle

# Note: in order to submit patches more easily, you might want to get a github
# account, fork the bot, and clone from your own writeable version.

# If you do that, clone with:
git clone git@github.com:<username>/sp0rkle.git
# and then add my repository as an alternative remote to pull from:
cd sp0rkle
git remote add -f -m master fluffle http://github.com/fluffle/sp0rkle.git

# Otherwise, just clone from my repository:
git clone http://github.com/fluffle/sp0rkle.git
```

5) Import the old databases into mongo:
 ```bash
go install github.com/fluffle/sp0rkle/util/importers/factimporter
factimporter --db=/path/to/db
go install github.com/fluffle/sp0rkle/util/importers/quoteimporter
quoteimporter --db=/path/to/db
```
# If you don't know where to get the dbs, you shouldn't be submitting patches :-)

6) Code, build, commit, push :)
 ```bash
while coding in $GOPATH/src/github.com/fluffle/sp0rkle/sp0rkle:
  vim <stuff>:wq
  go build
  # Run local build for testing ...
  ./sp0rkle --servers irc.pl0rt.org[:port]  [--nick=mybot] [--channels='#test']
  ^C

git add <stuff>
git commit -m "Some useful message about the edit to <stuff>."
```
# If you cloned from your own repo:
git push  # pushes changes in your branches up to github
# ... then send me a pull request on github :-)

# Otherwise, I guess you'll have to mail me a patch, or something:
# This might work, untested, you should read man git-format-patch(1).
git format-patch --attach --stdout --to=abramley@gmail.com | mail