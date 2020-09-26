package helpers

import (
	"encoding/json"
	"errors"
	"regexp"
)

type SecretFile struct {
	APIVersion string             `json:"apiVersion"`
	Data       interface{}        `json:"data"`
	Kind       string             `json:"kind"`
	Metadata   SecretFileMetaData `json:"metadata"`
	Type       string             `json:"type"`
}

type SecretFileMetaData struct {
	Name string `json:"name"`
}

func GenerateJSONSecret(name string, data interface{}) (string, error) {
	jsonValue, jsonError := json.Marshal(SecretFile{
		APIVersion: "v1",
		Data:       data,
		Kind:       "Secret",
		Metadata: SecretFileMetaData{
			Name: name,
		},
		Type: "Opaque",
	})
	return string(jsonValue), jsonError
}

var secretNameRegexFormat *regexp.Regexp = regexp.MustCompile("^" + "[a-z0-9]([-a-z0-9]*[a-z0-9])?" + "(\\." + "[a-z0-9]([-a-z0-9]*[a-z0-9])?" + ")*" + "$")

const SecretNameMaxLength int = 253
const InvalidSecretNameFormat string = "INVALID_SECRET_NAME_FORMAT"
const InvalidSecretNameLength string = "INVALID_SECRET_NAME_LENGTH"

func ValidateSecretName(name string) error {
	if len(name) > SecretNameMaxLength || len(name) == 0 {
		return errors.New(InvalidSecretNameLength)
	}
	if !secretNameRegexFormat.MatchString(name) {
		return errors.New(InvalidSecretNameFormat)
	}

	return nil
}
