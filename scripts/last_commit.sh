#!/bin/sh

if [ -z $1 ]; then
    echo '"$1" should be service-name' >&2
    exit 1
fi
export PKG=$1
cd $PKG
# Use go-list to dependencies of a package // HL
IMPORTS=$(go list -deps|grep 'github.com/laher/toasties-galore'|sed 's,github.com/laher/toasties-galore/,,g'|xargs)
cd ..
# Use git-rev-list to find last commit affecting dirs // HL
LAST_COMMIT=$(git rev-list -1 HEAD -- $IMPORTS $PKG ./go.mod ./go.sum)
echo $LAST_COMMIT
