#!/bin/zsh

# Vars
# Host
HOST="localhost"
PORT="8080"
PATH="routes"
# SRC
SRC_LAT="13.388860"
SRC_LNG="52.517037"
# DST 1
DST_LAT_1="13.397634"
DST_LNG_1="52.529407"
# DST 2
DST_LAT_2="13.428555"
DST_LNG_2="52.523219"

# Pre
# Curl and jq installed using nix not found if path not appropriately set
# Uncomment these helper lines or replace '/usr/bin/curl' by your system values
# if curl is not included in you PATH.
# curlcmd="$(which curl)"
# alias curl=$curlcmd

post () {
  echo "PUT $1"
  /usr/bin/curl -X GET $1
}


# Sample GET request: http://your-service/routes?src=13.388860,52.517037&dst=13.397634,52.529407&dst=13.428555,52.523219
post "http://$HOST:$PORT/$PATH?src=$SRC_LAT,$SRC_LNGi&$DST_LAT_1,$DST_LNG_1&$DST_LAT_2,$DST_LNG_2"

