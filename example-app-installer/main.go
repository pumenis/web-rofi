package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// === 1) SQL file URLs (raw) ===
	sqlURLs := []string{
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/ELZ.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/RST.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/OGB.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/NTPT.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/KJV.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/MNB.SQLite3.sql",
	}

	// === 2) Bash script URLs from cgi-bin dir ===
	scriptURLs := []string{
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/cgi-bin/modules.sh",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/cgi-bin/books.sh",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/cgi-bin/chapters.sh",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/cgi-bin/verses.sh",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/cgi-bin/display.sh",
		// add more scripts here
	}

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	dbTargetDir := filepath.Join(usr.HomeDir, ".local", "share", "mybible")
	if err := os.MkdirAll(dbTargetDir, 0o755); err != nil {
		panic(err)
	}

	for _, url := range sqlURLs {
		fmt.Println("Fetching SQL:", url)
		sqlText, err := fetchText(url)
		if err != nil {
			fmt.Println("Download error:", err)
			continue
		}

		fileName := path.Base(url)
		dbName := strings.TrimSuffix(fileName, ".sql")
		dbPath := filepath.Join(dbTargetDir, dbName)

		fmt.Println("Creating DB:", dbPath)
		if err := runSQLOnDB(dbPath, sqlText); err != nil {
			fmt.Println("DB error:", err)
		} else {
			fmt.Println("✔ DB created:", dbPath)
		}
	}

	cgiDir := filepath.Join(usr.HomeDir, ".cgi-bin")
	if err := os.MkdirAll(cgiDir, 0o755); err != nil {
		panic(err)
	}

	for _, url := range scriptURLs {
		fmt.Println("Fetching script:", url)
		if err := downloadFile(url, cgiDir, 0o755); err != nil {
			fmt.Println("Script error:", err)
		} else {
			fmt.Println("✔ Script saved to", cgiDir)
		}
	}
}

func fetchText(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad HTTP status: %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func runSQLOnDB(dbPath, sqlText string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec(sqlText)
	return err
}

func downloadFile(url, targetDir string, perm os.FileMode) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}
	// filename from URL path
	fileName := path.Base(url)
	targetPath := filepath.Join(targetDir, fileName)

	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}
	if err := os.Chmod(targetPath, perm); err != nil {
		return err
	}
	return nil
}
