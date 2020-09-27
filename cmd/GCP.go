package cmd

import (
	secretmanagerApiV1Beta1 "cloud.google.com/go/secretmanager/apiv1beta1"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/secretmanager/v1beta1"
	"log"
	"s5s/helpers"
	"strings"
	"sync"
)

type returnData struct {
	Key   string
	Value string
}

var gcpCommand = &cobra.Command{
	Use:   "gcp",
	Short: "GCP Secret Connector",
	Run: func(cmd *cobra.Command, args []string) {
		// K8s Flags
		k8sSecretName, _ := cmd.Flags().GetString("output-secret")
		k8sSecretNameError := helpers.ValidateSecretName(k8sSecretName)
		if k8sSecretNameError != nil {
			if k8sSecretNameError.Error() == helpers.InvalidSecretNameFormat {
				log.Fatal("K8s secret name must be a valid format")
			} else if k8sSecretNameError.Error() == helpers.InvalidSecretNameLength {
				log.Fatalf("K8s secret name must be a valid length of no greater than %d characters", helpers.SecretNameMaxLength)
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

		ch := make(chan returnData)
		var wg sync.WaitGroup

		for _, secret := range secrets {
			wg.Add(1)
			go func(secret string) {
				defer wg.Done()

				secretKV := strings.Split(secret, "=")
				if len(secretKV) != 2 {
					log.Fatal("Secret error: Secret must be in the \"k8sKey=gcpKey\" format")
				}

				k8sSecretKey := secretKV[0]
				gcpSecretKey := secretKV[1]
				secretName := "projects/" + project + "/secrets/" + gcpSecretKey + "/versions/" + version
				request := secretmanager.AccessSecretVersionRequest{
					Name: secretName,
				}

				if response, responseError := client.AccessSecretVersion(ctx, &request); responseError != nil {
					log.Fatal("Response error: " + responseError.Error())
				} else {
					ch <- returnData{
						Key:   k8sSecretKey,
						Value: string(response.Payload.Data),
					}
				}
			}(secret)
		}

		go func() {
			wg.Wait()
			close(ch)
		}()

		k8sSecrets := make(map[string]string)
		for data := range ch {
			k8sSecrets[data.Key] = data.Value
		}

		k8sSecretJSON, err := helpers.GenerateJSONSecret(k8sSecretName, k8sSecrets)
		if err != nil {
			log.Fatal("JSON Helper: " + err.Error())
		}

		fmt.Println(k8sSecretJSON)
	},
}

func init() {
	gcpCommand.Flags().StringP("key", "k", "", "GCP JSON Credentials as string (this or --key-file is required)")

	gcpCommand.Flags().StringP("key-file", "f", "", "GCP JSON File Credentials (this or --key is required)")

	gcpCommand.Flags().StringP("project", "p", "", "GCP Project Name (required)")
	mustMarkFlagRequired("project")

	gcpCommand.Flags().StringP("version", "v", "latest", "GCP Secret Version (default: latest)")

	gcpCommand.Flags().StringArrayP("secret", "s", []string{}, "Kubernetes Secret Name")
	mustMarkFlagRequired("secret")
}

func mustMarkFlagRequired(name string) {
	if markRequiredError := gcpCommand.MarkFlagRequired("project"); markRequiredError != nil {
		log.Fatal(markRequiredError)
	}
}
