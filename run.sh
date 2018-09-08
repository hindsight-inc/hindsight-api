#!/bin/bash
if go build main.go; then
	./main
else
	echo ---- Build failed ----
fi