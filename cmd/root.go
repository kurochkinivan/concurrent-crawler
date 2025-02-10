package cmd

import (
	"context"
	"log"
	"time"

	"github.com/kurochkinivan/web-local-mirror/internal/mirror"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().StringP("website", "w", "https://go.dev/", "Website to mirror, default: go.dev")
	rootCmd.Flags().Int("workers", 20, "Number of workers, default: 20")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var rootCmd = &cobra.Command{
	Use: "web-local-mirror",
	Run: func(cmd *cobra.Command, args []string) {
		website, _ := cmd.Flags().GetString("website")
		workerCount, _ := cmd.Flags().GetInt("workers")
		runMirror(website, workerCount)
	},
}

func runMirror(website string, workerCount int) {
	log.Println("starting mirroring the website")
	start := time.Now()
	m, err := mirror.NewMirror(website, workerCount)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("mirroring complete! Total time: %v\n", time.Since(start))
}
