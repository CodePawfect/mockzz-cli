package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/CodePawfect/mockzz-cli/internals"
	"github.com/CodePawfect/mockzz-cli/model"
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

		endpoints, err := model.ReadEndpoints("mockzz-endpoints.txt")
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

		port := cmd.Flag("port").Value.String()
		isValidPort(port)

		internals.StartHttpServer(mux, port)
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

func isValidPort(portStr string) {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Error converting port to integer: %v", err)
	}

	if port < 0 || port > 65535 {
		log.Fatalf("Port is not in range >= 0 && <= 65535: %v", err)
	}
}
