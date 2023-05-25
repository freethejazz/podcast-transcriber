package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
)

func downloadMP3(jobFolder string, url string) (string, error) {
	// Get the filename from the URL
	_, filename := path.Split(url)
	filename = path.Base(filename)

	// Create the output directory if it doesn't exist
	err := os.MkdirAll(path.Join("./dls", jobFolder), os.ModePerm)
	if err != nil {
		return "", err
	}

	// Create the output file
	filepath := path.Join("./dls", jobFolder, filename)
	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Send HTTP GET request to download the file
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Check if the response was successful
	if response.StatusCode != http.StatusOK {
		return "", err
	}

	// Copy the response body to the output file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", err
	}

	log.Printf(filename)

	return filename, nil
}
