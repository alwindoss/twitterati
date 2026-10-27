package engine

import (
	"errors"
	"fmt"
	"os"
	"time"

	"go.dalw.in/twitterati/internal/twitter"
)

func DownloadVideoFromTwitter(url string, outputFileName string) error {
	var file *os.File
	file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			newFileName := fmt.Sprintf("%s_%s", time.Now().String(), outputFileName)
			file, err = os.OpenFile(newFileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println("Error creating file:", err)
				return err
			}
		} else {
			return err
		}
	}
	defer file.Close()

	fmt.Println("File created successfully!")

	cfg := &twitter.Config{
		URL:          url,
		OutputWriter: file,
	}
	err = twitter.Run(cfg)
	if err != nil {
		return err
	}
	return nil
}
