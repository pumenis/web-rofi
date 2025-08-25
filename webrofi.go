package main

import "fmt"

func GetHTML(list []string, url string) string {
	options := ""

	for _, option := range list {
		options = options +
			`<option value="` +
			option +
			`"></option>` + "\n"
	}

	return fmt.Sprintf(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title></title>
  </head>
  <body>
    <input id="mainInput" list="l-0" onkeyup="handlekeyup(event)" class="dropdown" type="text"></input>
    <datalist id="l-0">
    		%s
    </datalist>  
    <script>
    	window.onload = function(){
				document.getElementById("mainInput").focus()
    	}
      function handlekeyup(e){
        if (e.key == undefined||e.key =="Enter") {
       	  window.location.assign("%s&arg=" + encodeURIComponent(e.target.value)+ '#bm') 
  			}	
      }
    </script>
  </body>
</html>`, options, url)
}
