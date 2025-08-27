#!/usr/bin/env bash

path="$HOME/.local/share/mybible/"

long_name=$(sed "s/'/''/g" <<<"$2")
echo "http://localhost:1121/rofi?script=verses.sh&arg=$1&arg=$2"
sqlite3 "${path}${1^^}.SQLite3" " -- sql
    SELECT v.chapter
      FROM verses v
INNER JOIN books b ON v.book_number = b.book_number
     WHERE b.long_name = '$long_name' AND v.verse = 1
"
