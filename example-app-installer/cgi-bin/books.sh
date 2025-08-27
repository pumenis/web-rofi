#!/usr/bin/env bash

path="$HOME/.local/share/mybible/"

echo "http://localhost:1121/rofi?script=chapters.sh&arg=$1"
sqlite3 "${path}${1^^}.SQLite3" " -- sql
  SELECT long_name
    FROM books
"
