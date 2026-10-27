package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.dalw.in/twitterati/internal/engine"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Downloads the video from a tweet given its URL",
	Long: `This command allows you to download the video from a tweet by providing the tweet's URL.
You can also specify an optional output filename. If not provided, the video will be saved with the tweet ID as the filename.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFileName == "" {
			outputFileName = fmt.Sprintf("output_%s.mp4", time.Now().String())
		}
		return engine.DownloadVideoFromTwitter(url, outputFileName)
	},
}

var (
	url            string
	outputFileName string
)

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	downloadCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the tweet to download")
	downloadCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "Output filename (optional)")
	downloadCmd.MarkFlagRequired("url")
}
