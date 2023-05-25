package main

import (
	"log"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	r.POST("/job", func(c *gin.Context) {
		// Parse JSON
		var json struct {
			Url string `json:"url" binding:"required"`
		}

		var jobId = uuid.New()

		if c.Bind(&json) == nil {
			go func() {
				// Download the mp3 file
				log.Printf("Downloading %s", json.Url)
				filename, _ := downloadMP3(jobId.String(), json.Url)
				folderPath := path.Join("./dls", jobId.String())
				log.Printf("Downloaded %s to %s", json.Url, folderPath)

				// Transcribe the mp3 file
				log.Printf("Transcribing %s", filename)
				Transcribe(folderPath, filename)
				log.Printf("Finished transcribing %s", filename)

				// Parse the SRT file
				log.Printf("Parsing raw SRT captions for %s", filename)
				rawCaptions, _ := ParseSRT(path.Join(folderPath, filename+".srt"))
				log.Printf("Raw captions are parsed")

				// Process raw captions to add context and a parent ID
				log.Printf("Processing captions to include context")
				captions := ProcessRawCaptions(json.Url, rawCaptions)
				log.Printf("Processed captions")

				// Index processed captions to Elasticsearch
				log.Printf("Indexing processed captions to elasticsearch")
				err := IndexCaptions(captions)
				if err != nil {
					log.Fatalf("Failed to index captions %v", err)
				}
				log.Printf("Indexed captions")
			}()

			c.JSON(http.StatusOK, gin.H{"status": "started", "jobId": jobId.String()})
		}

	})

	r.POST("/search", func(c *gin.Context) {
		// Parse JSON
		var json struct {
			Query string `json:"query" binding:"required"`
		}

		if c.Bind(&json) == nil {
			captions, _ := SearchCaptions(json.Query)

			c.JSON(http.StatusOK, gin.H{"results": captions})
		}
	})

	return r
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
