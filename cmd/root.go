package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use: "s5s",
	Short: "Secrets Manager Bridge for Kubernetes",
}

func init() {
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "k8s namespace (required)")
	if err := rootCmd.MarkPersistentFlagRequired("namespace"); err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

