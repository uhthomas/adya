#!/bin/sh

if [ -z $GIT_COMMIT ]; then
	GIT_COMMIT="$(git rev-parse head)"
fi

if [ -z $GIT_REF ]; then
	GIT_REF="$(git tag --points-at head)"
fi

cat <<EOF
STABLE_GIT_COMMIT $GIT_COMMIT
STABLE_GIT_REF $GIT_REF
EOF
