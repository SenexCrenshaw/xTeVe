#!/bin/bash
webUIHeader="package src\n\nvar webUI = make(map[string]interface{})\n\nfunc loadHTMLMap() {\n"

if [ -f "../src/webUI.go" ]; then
    rm "../src/webUI.go"
fi

echo -e $webUIHeader > '../src/webUI.go'

for file in ../html/* ../html/**/*
do
    if [[ -f $file ]]; then
        filePath=$file
	fileName="${file/..\/}"
	base=$( base64 -w 0 "$filePath" )
	output="\twebUI[\"$fileName\"] = \"$base\""
	echo -e "$output" >> '../src/webUI.go'
    fi
done

echo -e "\n}" >> '../src/webUI.go'
