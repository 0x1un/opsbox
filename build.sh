#!/usr/bin/env bash


GO=$(which go)

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "go compiler not found"
  exit 127
fi

$GO mod download

# shellcheck disable=SC2181
if [ $? -ne 0 ]; then
  echo "dependencies download failed"
  exit 127
fi


rootDir=$(pwd)

declare -A map
map[feishu]='alertmanager/feishu'
map[gomi]='cli/linux/gomi'
map[bigdata-exporter]='exporter/bigdata'
map[chain-listener]='exporter/chain-listener/cmd'


function BuildProject() {
  output="$1"
  platform="$2"
  project="${map[$1]}"
  if [ "$platform" == "windows" ]; then
    cd "$project" && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $GO build -ldflags "-s -w" -o "${rootDir}/build/${output}.exe"
  else
    cd "$project" && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $GO build -ldflags "-s -w" -o "${rootDir}/build/${output}"
  fi
}

BuildProject "$1" "$2"