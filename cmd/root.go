package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "s5s",
	Short: "Secrets Manager Bridge for Kubernetes",
}

func init() {
	rootCmd.PersistentFlags().StringP("output-secret", "o", "", "k8s secret name (required)")
	_ = rootCmd.MarkPersistentFlagRequired("output-secret")
	rootCmd.AddCommand(gcpCommand)
	rootCmd.AddCommand(awsCommand)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
