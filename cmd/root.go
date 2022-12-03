package cmd

import (
	"github.com/Twacqwq/godown/download"
	"github.com/spf13/cobra"
)

var (
	url         string
	output      string
	concurrency int
)

var root = &cobra.Command{
	Use:   "godown",
	Short: "A Golang-based command line concurrent download tool",
	Run: func(cmd *cobra.Command, args []string) {
		download.NewGodown(concurrency, output).Download(url)
	},
	Args: cobra.NoArgs,
}

func init() {
	root.Flags().StringVarP(&url, "url", "u", "", "specify the download link address")
	root.Flags().StringVarP(&output, "output", "o", "", "select the download path")
	root.Flags().IntVarP(&concurrency, "concurrency", "n", 8, "specify concurrency")
}

func Execute() error {
	return root.Execute()
}
