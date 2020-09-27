package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/spf13/cobra"
	"log"
	"s5s/helpers"
)

var awsCommand = &cobra.Command{
	Use:   "aws",
	Short: "AWS Secret Connector",
	Run: func(cmd *cobra.Command, args []string) {
		// K8S Flags
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

		// AWS Flags
		awsSecretName, _ := cmd.Flags().GetString("secret")
		awsRegion, _ := cmd.Flags().GetString("region")
		awsAccessKeyId, _ := cmd.Flags().GetString("accessKeyId")
		awsAccessKey, _ := cmd.Flags().GetString("accessKey")
		awsToken, _ := cmd.Flags().GetString("token")

		config := aws.NewConfig()
		config.WithRegion(awsRegion)
		if awsAccessKeyId != "" && awsAccessKey != "" {
			creds := credentials.NewStaticCredentials(awsAccessKeyId, awsAccessKey, awsToken)
			config.WithCredentials(creds)
		} else {
			config.WithCredentials(credentials.NewEnvCredentials())
		}

		awsSession, sessionError := session.NewSession(config)
		if sessionError != nil {
			log.Fatal(sessionError)
		}
		svc := secretsmanager.New(awsSession)

		input := &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(awsSecretName),
		}

		secretValue, getSecretValueError := svc.GetSecretValue(input)
		if getSecretValueError != nil {
			if awsErr, ok := getSecretValueError.(awserr.Error); ok {
				log.Fatal(awsErr)
			} else {
				log.Fatal(getSecretValueError)
			}
		}

		awsSecretResponse := make(map[string]interface{})
		if jsonError := json.Unmarshal([]byte(*secretValue.SecretString), &awsSecretResponse); jsonError != nil {
			log.Fatal(jsonError)
		}

		k8sSecretJSON, generationError := helpers.GenerateJSONSecret(k8sSecretName, awsSecretResponse["data"])
		if generationError != nil {
			log.Fatal(generationError)
		}

		fmt.Println(k8sSecretJSON)
	},
}

func init() {
	awsCommand.Flags().String("accessKeyId", "", "AWS Access Key")
	awsCommand.Flags().String("accessKey", "", "AWS Secret Key")
	awsCommand.Flags().String("token", "", "AWS Session Token")

	awsCommand.Flags().StringP("region", "r", "us-west-2", "AWS Region Name (default: us-west-2)")

	awsCommand.Flags().StringP("secret", "s", "", "AWS Secret Name")
	if markRequiredError := awsCommand.MarkFlagRequired("secret"); markRequiredError != nil {
		log.Fatal(markRequiredError)
	}
}
