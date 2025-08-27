#!/usr/bin/env bash

path="~/.local/share/mybible/"

long_name=$(sed "s/'/''/g" <<<"$2")
echo '<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title></title>
  </head>
  <body>'

sqlite3 "${path}${1^^}.SQLite3" " -- sql
    SELECT
      CASE
        WHEN v.chapter = '$3' AND v.verse = '$4'
          THEN '<span id=\"bm\">$4</span> ' || v.text ||
            '<br>'
        ELSE v.verse || ' ' || v.text ||'<br>'
      END
      FROM verses v
INNER JOIN books b ON v.book_number = b.book_number
     WHERE b.long_name = '$long_name' 
"

echo "<script>
window.onload = function() {
  const bm = document.getElementById('bm');
  if (bm) bm.scrollIntoView();
};
</script>
</body>
</html>"
