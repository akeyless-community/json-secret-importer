package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

type Secret struct {
	Date       string `json:"date"`
	Secret     string `json:"secret"`
	Status     string `json:"status"`
	ValidAfter string `json:"valid_after"`
}

func main() {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
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

	for _, version := range activeVersions {
		decodedSecret, err := base64.StdEncoding.DecodeString(secrets[strconv.Itoa(version)].Secret)
		if err != nil {
			fmt.Println("Error decoding secret:", err)
			return
		}
		fmt.Printf("Version: %d, Secret: %s\n", version, string(decodedSecret))
	}
}
