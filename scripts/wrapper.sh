#!/bin/sh

if [ "$1" = "new" ]; then
	hugo new "post/$(date +'%Y-%m-%d')-$2.md"
else
	hugo $*
fi
