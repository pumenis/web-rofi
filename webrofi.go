package main

import "fmt"

func GetHTML(list []string, url string) string {
	options := ""

	for _, option := range list {
		options = options +
			"<div>" + option + "</div>\n"
	}

	return fmt.Sprintln(`
<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Rofi-like Dropdown with Clear Button</title>
    <style>
      .dropdown {
        position: fixed;
        top: 20%;
        left: 50%;
        transform: translateX(-50%);
        width: 300px;
        background: white;
        box-shadow: 0 4px 10px rgba(0, 0, 0, 0.2);
        padding: 10px;
        border-radius: 6px;
        font-family: sans-serif;
      }
      .input-wrapper {
        position: relative;
        width: 100%;
      }
      #myInput {
        width: 100%;
        padding: 10px 28px 10px 10px; /* space for clear button */
        font-size: 16px;
        box-sizing: border-box;
        border: 2px solid #4da3ff;
        outline: none;
      }
      .clear-btn {
        position: absolute;
        right: 8px;
        top: 50%;
        transform: translateY(-50%);
        cursor: pointer;
        font-weight: bold;
        color: #999;
        user-select: none;
        display: none; /* show only when input has text */
      }
      .clear-btn:hover {
        color: #333;
      }
      .dropdown-list {
        max-height: 300px;
        overflow-y: auto;
        border: 1px solid #ccc;
        background: white;
        margin-top: 5px;
      }
      .dropdown-list div {
        padding: 10px;
        cursor: pointer;
      }
      .dropdown-list div:hover,
      .dropdown-list .highlighted {
        background-color: #bde4ff;
      }
    </style>
  </head>
  <body>
    <div id="myDropdown" class="dropdown">
      <div class="input-wrapper">
        <input
          id="myInput"
          type="text"
          placeholder="Type to filter..."
          autocomplete="off"
        />
        <span id="clearBtn" class="clear-btn">&times;</span>
      </div>
      <div id="myList" class="dropdown-list">
    `) + fmt.Sprintf(`
    		%s
      </div>
    </div>

    <script>
      const input = document.getElementById("myInput");
      const list = document.getElementById("myList");
      const clearBtn = document.getElementById("clearBtn");
      let highlightedIndex = -1;

      // Always focus
      window.addEventListener("load", () => input.focus());
      input.addEventListener("blur", () => setTimeout(() => input.focus(), 0));

      function selectItem(value) {
        input.value = value;
        filterList(input.value);
				window.location.assign("%s&arg=" + encodeURIComponent(value)+ '#bm')
      }

      function highlightItem(index) {
        const children = [...list.querySelectorAll("div")];
        if (highlightedIndex >= 0 && highlightedIndex < children.length) {
          children[highlightedIndex].classList.remove("highlighted");
        }
        if (index >= 0 && index < children.length) {
          children[index].classList.add("highlighted");
          highlightedIndex = index;
        } else {
          highlightedIndex = -1;
        }
      }

      function attachItemListeners() {
        const children = list.querySelectorAll("div");
        children.forEach((child, idx) => {
          child.addEventListener("click", () => selectItem(child.textContent));
          child.addEventListener("mouseenter", () => highlightItem(idx));
        });
      }

      function filterList(filterValue) {
        const filter = filterValue.trim().toLowerCase();
        const children = list.querySelectorAll("div");
        children.forEach((child) => {
          child.style.display = child.textContent.toLowerCase().includes(filter)
            ? "block"
            : "none";
        });
        highlightedIndex = -1;
        clearBtn.style.display = filterValue ? "block" : "none";
      }

      input.addEventListener("input", () => filterList(input.value));

      clearBtn.addEventListener("click", () => {
        input.value = "";
        filterList("");
        input.focus();
      });

      input.addEventListener("keydown", (e) => {
        const visibleItems = [...list.querySelectorAll("div")].filter(
          (c) => c.style.display !== "none",
        );
        if (e.key === "ArrowDown") {
          e.preventDefault();
          if (highlightedIndex < visibleItems.length - 1) {
            highlightItem(highlightedIndex + 1);
            visibleItems[highlightedIndex].scrollIntoView({ block: "nearest" });
          }
        } else if (e.key === "ArrowUp") {
          e.preventDefault();
          if (highlightedIndex > 0) {
            highlightItem(highlightedIndex - 1);
            visibleItems[highlightedIndex].scrollIntoView({ block: "nearest" });
          }
        } else if (e.key === "Enter") {
          e.preventDefault();
          if (highlightedIndex >= 0) {
            selectItem(visibleItems[highlightedIndex].textContent);
          } else if (visibleItems.length === 1) {
            selectItem(visibleItems[0].textContent);
          }
        }
      });

      attachItemListeners();
    </script>
  </body>
</html>
`, options, url)
}
