#!/bin/bash
rm hindsight
echo ---- Building ----
if go build; then
	echo ---- Testing ----
	if go test; then
		echo ---- Running ----
		./hindsight
	else
		echo ---- Test failed ----
	fi
else
	echo ---- Build failed ----
fi