package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/CodePawfect/mockzz-cli/internals"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const logo = `

 __  __  ___   _______ ______________
|  \/  |/ _ \ / ___| |/ /___  /___  /
| \  / | | | | |   | ' /   / /   / /
| |\/| | | | | |   |  <   / /   / / 
| |  | | |_| | |___| . \ / /__ / /__
|_|  |_|\___/ \____|_|\_\_____|_____|

`

var logoStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#01FAC6")).Bold(true)

var port string

func init() {
	startCmd.PersistentFlags().StringVarP(&port, "port", "p", "8080", "This flag sets the port of our API server")
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the HTTP Server to serve the mocked APIs",
	Long:  `This command starts the HTTP server that serves the mocked APIs. It listens on the specified port and serves the mock responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", logoStyle.Render(logo))

		endpoints, err := readEndpoints("mockzz-endpoints.txt")
		if err != nil {
			fmt.Println("Error reading endpoints:", err)
			return
		}

		if len(endpoints) > 0 {
			fmt.Println("Mocked APIs:")
			fmt.Println()
		}

		for api := range endpoints {
			fmt.Printf("API: %s\n", api)
		}

		mux := http.NewServeMux()

		for api, responseFile := range endpoints {
			handlerFunc := createHandlerFunc(responseFile)
			mux.HandleFunc(api, handlerFunc)
		}

		internals.StartHttpServer(mux, cmd.Flag("port").Value.String())
	},
}

func createHandlerFunc(responseFilePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := os.ReadFile(responseFilePath)
		if err != nil {
			fmt.Printf("Error reading response file %s: %v\n", responseFilePath, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}

// readEndpoints reads the endpoints from the given file and returns a map[api]responseFile
func readEndpoints(filename string) (map[string]string, error) {
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
