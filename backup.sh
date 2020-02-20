#! /bin/sh

set -e

BACKUPDIR="$1"
if [ ! -d "$BACKUPDIR" ]; then
    echo "Backup path '$BACKUPDIR' is not a directory :-(" >&2
    exit 1
fi

TMPDIR="$(/bin/mktemp -d)"
trap "/bin/rm -r \"$TMPDIR\"; exit" INT TERM EXIT
DATE="$(date "+%Y-%m-%d.%H:%M")"
/usr/bin/mongodump -d sp0rkle -o "$TMPDIR" >/dev/null
/bin/tar -C "$TMPDIR" -cjf "$BACKUPDIR/sp0rkle.$DATE.tar.bz2" .

# To restore into bitnami mongodb docker image in k8s:
#  - untar to persistent volume so container can see backup
#  - docker exec -t <container> /opt/bitnami/mongodb/bin/mongorestore \
#        --drop -d sp0rkle --dir=/bitnami/mongodb/restore/sp0rkle -v
