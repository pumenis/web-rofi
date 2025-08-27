package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Direct RAW URLs to SQL files from your GitHub repo
	sqlURLs := []string{
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/ELZ.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/RST.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/OGB.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/NTPT.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/KJV.SQLite3.sql",
		"https://github.com/pumenis/web-rofi/raw/refs/heads/main/example-app-installer/modules/MNB.SQLite3.sql",
	}

	// Find home dir
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	targetDir := filepath.Join(usr.HomeDir, ".local", "share", "mybible")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		panic(err)
	}

	for _, url := range sqlURLs {
		fmt.Println("Fetching:", url)

		sqlText, err := fetchSQL(url)
		if err != nil {
			fmt.Println("Download error:", err)
			continue
		}

		_, file := filepath.Split(url)
		dbName := strings.TrimSuffix(file, ".sql") + ".db"
		dbPath := filepath.Join(targetDir, dbName)

		fmt.Println("Creating DB:", dbPath)
		if err := runSQLOnDB(dbPath, sqlText); err != nil {
			fmt.Println("DB error:", err)
		} else {
			fmt.Println("✔ Done:", dbPath)
		}
	}
}

func fetchSQL(url string) (string, error) {
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
