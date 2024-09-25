package model

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Endpoint struct {
	API          string
	ResponseFile string
}

type Endpoints struct {
	Endpoints []Endpoint
}

func NewModel() *Endpoints {
	return &Endpoints{}
}

func LoadModel() *Endpoints {
	responseByApi, err := ReadEndpoints("mockzz-endpoints.txt")
	if err != nil {
		fmt.Println("Error reading endpoints:", err)
		return NewModel()
	}

	endpoints := make([]Endpoint, 0)

	for api, responseFile := range responseByApi {
		e := Endpoint{
			API:          api,
			ResponseFile: responseFile,
		}

		endpoints = append(endpoints, e)
	}

	return &Endpoints{endpoints}
}

func (m *Endpoints) GetEndpoints() []Endpoint {
	return m.Endpoints
}

func (m *Endpoints) AddEndpoint(e Endpoint) {
	m.Endpoints = append(m.Endpoints, e)
}

// readEndpoints reads the endpoints from the given file and returns a map[api]responseFilePath
func ReadEndpoints(filename string) (map[string]string, error) {
	endpoints := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines and lines without a colon
		if strings.TrimSpace(line) == "" || !strings.Contains(line, ":") {
			continue
		}

		// Split the line into API and response file
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			fmt.Printf("Invalid line format (skipping): %s\n", line)
			continue
		}
		api := strings.TrimSpace(parts[0])
		responseFile := strings.TrimSpace(parts[1])

		// Add to the map
		endpoints[api] = responseFile
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return endpoints, nil
}
