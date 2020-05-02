package cmd

import (
	secretmanagerApiV1Beta1 "cloud.google.com/go/secretmanager/apiv1beta1"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"log"
	"os"
	"s5s/helpers"
	"strings"
)

func init() {
	gcpCommand.Flags().StringP("key", "k", "", "GCP JSON Credentials as string (this or --key-file is required)")

	gcpCommand.Flags().StringP("key-file", "f", "", "GCP JSON File Credentials (this or --key is required)")

	gcpCommand.Flags().StringP("project", "p", "", "GCP Project Name (required)")
	_ = gcpCommand.MarkFlagRequired("project")

	gcpCommand.Flags().StringP("version", "v", "latest", "GCP Secret Version (default: latest)")

	gcpCommand.Flags().StringArrayP("secret", "s", []string{}, "Secrets")
	_ = gcpCommand.MarkFlagRequired("secret")

	rootCmd.AddCommand(gcpCommand)
}

var gcpCommand = &cobra.Command{
	Use:   "gcp",
	Short: "GCP Secret Connector",
	Run: func(cmd *cobra.Command, args []string) {
		// K8s Flags
		k8sSecretName, _ := cmd.Flags().GetString("output-secret")
		k8sSecretNameError := helpers.ValidateSecretName(k8sSecretName)
		if len(k8sSecretNameError) > 0 {
			if k8sSecretNameError == helpers.InvalidSecretNameFormat {
				fmt.Println("K8s secret name must be a valid format")
				os.Exit(1)
			} else if k8sSecretNameError == helpers.InvalidSecretNameLength {
				fmt.Println(fmt.Sprintf(
					"K8s secret name must be a valid length of no greater than %d characters",
					helpers.SecretNameMaxLength))
				os.Exit(1)
			}
			log.Fatal(k8sSecretNameError)
		}

		// GCP Flags
		key, _ := cmd.Flags().GetString("key")
		keyFile, _ := cmd.Flags().GetString("key-file")
		project, _ := cmd.Flags().GetString("project")
		secrets, _ := cmd.Flags().GetStringArray("secret")
		version, _ := cmd.Flags().GetString("version")

		// Mark required
		if key == "" && keyFile == "" {
			_ = cmd.MarkFlagRequired("key")
			_ = cmd.MarkFlagRequired("key-file")
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

		k8sSecrets := make(map[string]string)

		channel := make(chan func() (string, string))

		for _, secret := range secrets {
			go func(secret string) {
				secretKV := strings.Split(secret, "=")
				if len(secretKV) != 2 {
					log.Fatal("Secret error: Secret must be in the \"k8sKey=gcpKey\" format")
				}

				k8sSecretKey := secretKV[0]
				gcpSecretKey := secretKV[1]
				request := secretmanager.AccessSecretVersionRequest{
					Name:                 "projects/" + project + "/secrets/" + gcpSecretKey + "/versions/" + version,
					XXX_NoUnkeyedLiteral: struct{}{},
					XXX_unrecognized:     nil,
					XXX_sizecache:        0,
				}

				if response, responseError := client.AccessSecretVersion(ctx, &request); responseError != nil {
					log.Fatal("Response error: " + responseError.Error())
				} else {
					returnValues := func(c chan func() (string, string)) {
						c <- func() (string, string) {
							return k8sSecretKey, base64.StdEncoding.EncodeToString(response.Payload.Data)
						}
					}
					go returnValues(channel)
				}
			}(secret)
		}

		secretsFound := 0
		for v := range channel {
			secretsFound++
			k8sKey, secretValue := v()
			k8sSecrets[k8sKey] = secretValue
			if secretsFound == len(secrets) {
				close(channel)
			}
		}

		k8sSecret, err := helpers.GenerateJSONSecret(k8sSecretName, k8sSecrets)
		if err != nil {
			log.Fatal("JSON Helper: " + err.Error())
		}

		fmt.Println(k8sSecret)
	},
}
