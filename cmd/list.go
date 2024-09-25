package cmd

import (
	"fmt"
	"os"

	"github.com/CodePawfect/mockzz-cli/model"
	ui "github.com/CodePawfect/mockzz-cli/ui/teaList"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows an interactive list of mocked APIs",
	Long:  `This command shows an interactive list of all the mocked APIs that have been created. It allows the user to view, edit, and delete the APIs.`,
	Run: func(cmd *cobra.Command, args []string) {
		dataModel := model.LoadModel()

		uiModel := ui.NewModel(dataModel)

		p := tea.NewProgram(uiModel)
		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	},
}
