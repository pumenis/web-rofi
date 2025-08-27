#!/usr/bin/env bash

path="$HOME/.local/share/mybible/"

long_name=$(sed "s/'/''/g" <<<"$2")
echo "http://localhost:1121/view?script=display.sh&arg=$1&arg=$2&arg=$3"
sqlite3 "${path}${1^^}.SQLite3" " -- sql
    SELECT v.verse
      FROM verses v
INNER JOIN books b ON v.book_number = b.book_number
     WHERE b.long_name = '$long_name' AND v.chapter = '$3'
"
