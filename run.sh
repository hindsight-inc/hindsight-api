#!/bin/bash
if go build ./...; then
	./main
else
	echo ---- Build failed ----
fi