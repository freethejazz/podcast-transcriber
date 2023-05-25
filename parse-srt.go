package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// RawCaption struct represents a caption with text, timestamps, and clip length.
type RawCaption struct {
	Index         int
	Text          string
	TimestampFrom time.Duration
	TimestampTo   time.Duration
	ClipLength    time.Duration
}

// ParseSRT parses an SRT file into a list of RawCaption structs.
func ParseSRT(filename string) ([]RawCaption, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var captions []RawCaption
	var currentCaption RawCaption
	var captionText strings.Builder
	var timestampRegex = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})\s+-->\s+(\d{2}):(\d{2}):(\d{2}),(\d{3})`)
	var indexRegex = regexp.MustCompile(`^\d+$`)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Empty line indicates the end of a caption
			currentCaption.Text = captionText.String()
			captions = append(captions, currentCaption)

			// Reset variables for the next caption
			currentCaption = RawCaption{}
			captionText.Reset()
		} else if indexRegex.MatchString(line) {
			captionIndex := scanner.Text()
			currentCaption.Index, _ = strconv.Atoi(captionIndex)
		} else if timestampRegex.MatchString(line) {
			matches := timestampRegex.FindStringSubmatch(line)

			hoursFrom, _ := strconv.Atoi(matches[1])
			minutesFrom, _ := strconv.Atoi(matches[2])
			secondsFrom, _ := strconv.Atoi(matches[3])
			millisecondsFrom, _ := strconv.Atoi(matches[4])
			currentCaption.TimestampFrom = time.Duration(hoursFrom)*time.Hour +
				time.Duration(minutesFrom)*time.Minute +
				time.Duration(secondsFrom)*time.Second +
				time.Duration(millisecondsFrom)*time.Millisecond

			hoursTo, _ := strconv.Atoi(matches[5])
			minutesTo, _ := strconv.Atoi(matches[6])
			secondsTo, _ := strconv.Atoi(matches[7])
			millisecondsTo, _ := strconv.Atoi(matches[8])
			currentCaption.TimestampTo = time.Duration(hoursTo)*time.Hour +
				time.Duration(minutesTo)*time.Minute +
				time.Duration(secondsTo)*time.Second +
				time.Duration(millisecondsTo)*time.Millisecond

			currentCaption.ClipLength = currentCaption.TimestampTo - currentCaption.TimestampFrom
		} else {
			captionText.WriteString(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return captions, nil
}
