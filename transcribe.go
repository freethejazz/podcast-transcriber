package main

import (
	"bytes"
	"log"
	"os/exec"
	"path"
)

func Transcribe(folder string, filename string) error {
	cmd := exec.Command("whisper", path.Join(folder, filename), "--model", "tiny.en", "--output_dir", folder)

	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
