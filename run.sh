#!/bin/bash
rm hindsight
if go build; then
	./hindsight
else
	echo ---- Build failed ----
fi