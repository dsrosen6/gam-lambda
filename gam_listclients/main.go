package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type result struct {
	Clients []string `json:"clients"`
}

func getClients() (*result, error) {
	dirPath := "/opt/GAMConfig/clients"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, fmt.Errorf("ReadDir: %w", err)
	}

	var clients []string
	for _, file := range files {
		if file.IsDir() {
			clients = append(clients, file.Name())
		}
	}

	return &result{Clients: clients}, nil
}

func main() {
	lambda.Start(getClients)
}
