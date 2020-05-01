package helpers

import "encoding/json"

type SecretFile struct {
	APIVersion string `json:"apiVersion"`
	Data interface{} `json:"data"`
	Kind string `json:"kind"`
	Metadata SecretFileMetaData `json:"metadata"`
	Type string `json:"type"`
}

type SecretFileMetaData struct {
	Name string `json:"name"`
}

func GenerateJSONSecret(name string, data interface{}) ([]byte, error) {
	return json.Marshal(SecretFile{
		APIVersion: "v1",
		Data:       data,
		Kind:       "Secret",
		Metadata: 	SecretFileMetaData{
			Name: name,
		},
		Type: "Opaque",
	})
}