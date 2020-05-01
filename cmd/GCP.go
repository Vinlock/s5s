package cmd

import (
	secretmanagerApiV1Beta1 "cloud.google.com/go/secretmanager/apiv1beta1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"log"
)

func init() {
	gcpCommand.Flags().String("key", "", "JSON Credentials as string")
	gcpCommand.Flags().String("keyFile", "", "JSON File Credentials")
	gcpCommand.Flags().StringP("project", "p", "", "JSON File Credentials")
	gcpCommand.Flags().StringP("version", "v", "latest", "JSON File Credentials")
	_ = gcpCommand.MarkFlagRequired("project")

	rootCmd.AddCommand(gcpCommand)
}

var gcpCommand = &cobra.Command{
	Use: "gcp",
	Args: cobra.MinimumNArgs(1),
	Short: "GCP Secret Connector",
	Run: func(cmd *cobra.Command, args []string) {
		// K8s Flags
		k8sNamespace, _ := cmd.Flags().GetString("namespace")

		// GCP Flags
		key, _ := cmd.Flags().GetString("key")
		keyFile, _ := cmd.Flags().GetString("keyFile")
		project, _ := cmd.Flags().GetString("project")
		version, _ := cmd.Flags().GetString("version")

		// Mark required
		if key == "" && keyFile == "" {
			_ = cmd.MarkFlagRequired("key")
			_ = cmd.MarkFlagRequired("keyFile")
		}

		var credentials option.ClientOption
		if key != "" {
			credentials = option.WithCredentialsJSON([]byte(key))
		} else if keyFile != "" {
			credentials = option.WithCredentialsFile(keyFile)
		}

		ctx := context.Background()
		client, clientError := secretmanagerApiV1Beta1.NewClient(ctx, credentials)
		if clientError != nil {
			log.Fatal("Client error: " + clientError.Error())
		}

		request := secretmanager.AccessSecretVersionRequest{
			Name:                 "projects/" + project + "/secrets/" + args[0] + "/versions/" + version,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		}

		if response, responseError := client.AccessSecretVersion(ctx, &request); responseError != nil {
			log.Fatal("Response error: " + responseError.Error())
		} else {
			fmt.Println(k8sNamespace, response.Payload.Data)
		}
	},
}

