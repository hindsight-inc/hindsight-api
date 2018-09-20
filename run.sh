#!/bin/bash
rm -f hindsight
echo ---- Building ----
if go build; then
	echo ---- Testing ----
	if go test -p 1; then
		echo ---- Running ----
		./hindsight
	else
		echo ---- Test failed ----
	fi
else
	echo ---- Build failed ----
fi