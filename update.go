package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func updateDB() {
	serverLog("Updating database...\n")

	// Build the maxmind url for downloading
	dbURL := "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=" + Settings.License + "&suffix=tar.gz"

	// Download the new database
	err := downloadFile("maxmind/GeoLite2-City.tar.gz", dbURL)
	if err != nil {
		serverLog("error updating database [5]\n")
		run = false
		return
	}

	// Extracting the downloaded file
	serverLog("Extracting database...\n")
	file, err := os.Open("maxmind/GeoLite2-City.tar.gz")

	gzr, err := gzip.NewReader(file)
	if err != nil {
		serverLog("error updating database [4]\n")
		run = false
		return
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	waiting := true
	for waiting {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
		case err != nil:
			serverLog("error updating database [3]\n")
			run = false
			return
		case header == nil:
			continue
		}

		fileName := filepath.Base(header.Name)

		if header.Typeflag == tar.TypeReg && fileName == "GeoLite2-City.mmdb" {
			serverLog("Found database in archive file\n")

			f, err := os.OpenFile(filepath.Join("maxmind/", fileName), os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				serverLog("error updating database [1]\n")
				run = false
				return
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				serverLog("error updating database [2]\n")
				run = false
				return
			}

			f.Close()

			waiting = false
		}
	}

	os.Remove("maxmind/GeoLite2-City.tar.gz")

	serverLog("Archive File Deleted\n")
}

func downloadFile(filepath string, url string) error {

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
