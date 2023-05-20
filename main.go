package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	akeyless "github.com/akeylesslabs/akeyless-go/v2"
)

type Secret struct {
	Date       string `json:"date"`
	Secret     string `json:"secret"`
	Status     string `json:"status"`
	ValidAfter string `json:"valid_after"`
}

var token string
var importPath string

func main() {

	// Get AKEYLESS_TOKEN from environment variable or user input
	token = os.Getenv("AKEYLESS_TOKEN")
	if len(token) == 0 {
		fmt.Print("Enter AKEYLESS_TOKEN: ")
		fmt.Scanln(&token)
		if len(token) == 0 {
			fmt.Println("AKEYLESS_TOKEN is required")
			return
		}
	}

	// Get AKEYLESS_IMPORT_STARTING_PATH from environment variable or user input
	importPath = os.Getenv("AKEYLESS_IMPORT_STARTING_PATH")
	if len(token) == 0 {
		importPath = "."
	}

	err := filepath.Walk(importPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".json" {
			processFile(path)
		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func processFile(path string) {
	fmt.Println("Processing file:", path)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var secrets map[string]Secret
	err = json.Unmarshal(data, &secrets)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	var activeVersions []int
	for version, secret := range secrets {
		if secret.Status == "PRIMARY" || secret.Status == "ACTIVE" {
			v, err := strconv.Atoi(version)
			if err != nil {
				fmt.Println("Error parsing version:", err)
				return
			}
			activeVersions = append(activeVersions, v)
		}
	}

	sort.Ints(activeVersions)

	// Get AKEYLESS_SECRET_NAME_PREFIX from environment variable
	secretNamePrefix := os.Getenv("AKEYLESS_SECRET_NAME_PREFIX")

	// Get AKEYLESS_API_GW_URL from environment variable or use default
	apiGatewayURL := os.Getenv("AKEYLESS_API_GW_URL")
	if len(apiGatewayURL) == 0 {
		apiGatewayURL = "https://api.akeyless.io"
	}

	// Initialize Akeyless client

	client := akeyless.NewAPIClient(&akeyless.Configuration{
		Servers: []akeyless.ServerConfiguration{
			{
				URL: apiGatewayURL,
			},
		},
	}).V2Api

	// Create a static secret with the oldest version
	decodedSecret, err := base64.StdEncoding.DecodeString(secrets[strconv.Itoa(activeVersions[0])].Secret)
	if err != nil {
		fmt.Println("Error decoding secret:", err)
		return
	}

	secretName := secretNamePrefix + strings.TrimSuffix(path, ".json")

	selectedType := "generic"

	createSecretBody := akeyless.CreateSecret{
		Name:  secretName,
		Value: string(decodedSecret),
		Type:  &selectedType,
		Token: &token,
	}

	_, _, err = client.CreateSecret(context.Background()).Body(createSecretBody).Execute()
	if err != nil {
		fmt.Println("Error creating secret:", err)
		return
	}

	fmt.Printf("Created secret with name: %s and version: %d\n", secretName, activeVersions[0])

	// Update the secret with all other versions
	for _, version := range activeVersions[1:] {
		decodedSecret, err := base64.StdEncoding.DecodeString(secrets[strconv.Itoa(version)].Secret)
		if err != nil {
			fmt.Println("Error decoding secret:", err)
			return
		}

		keepPreviousVersion := "true"

		updateSecretBody := akeyless.UpdateSecretVal{
			Name:            secretName,
			Value:           string(decodedSecret),
			KeepPrevVersion: &keepPreviousVersion,
			Token:           &token,
		}

		_, _, err = client.UpdateSecretVal(context.Background()).Body(updateSecretBody).Execute()
		if err != nil {
			fmt.Println("Error updating secret:", err)
			return
		}

		fmt.Printf("Updated secret with name: %s and version: %d\n", secretName, version)
	}
}
