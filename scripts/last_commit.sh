#!/bin/sh
if [ -z $1 ]; then
    echo '"$1" should be service-name' >&2
    exit 1
fi
export PKG=$1
cd $PKG
IMPORTS=$(go list -deps|grep 'github.com/laher/toasties-galore'|sed 's,github.com/laher/toasties-galore/,,g'|xargs)
cd ..
LAST_COMMIT=$(git rev-list -1 HEAD -- $IMPORTS $PKG)
echo $LAST_COMMIT
