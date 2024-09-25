package ui

import (
	"fmt"
	"os"

	"github.com/CodePawfect/mockzz-cli/model"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	list      list.Model
	endpoints *model.Endpoints
	adding    bool
	apiInput  textinput.Model
	fileInput textinput.Model
}

func NewModel(m *model.Endpoints) Model {
	//fmt.Printf("%s\n", logoStyle.Render(logo))

	items := make([]list.Item, len(m.Endpoints))

	// Populate items with list.Items created from m.Endpoints
	for i, endpoint := range m.Endpoints {
		items[i] = listItem(endpoint)
	}

	delegate := newItemDelegate()
	l := list.New(items, delegate, 80, 14)
	l.Title = "Your Mocked Endpoints"
	l.SetShowStatusBar(false)

	// Initialize text inputs
	apiInput := textinput.New()
	apiInput.Placeholder = "/your/api/endpoint"
	apiInput.CharLimit = 256
	apiInput.Width = 50

	fileInput := textinput.New()
	fileInput.Placeholder = "path/to/response/file"
	fileInput.CharLimit = 256
	fileInput.Width = 50

	return Model{
		list:      l,
		endpoints: m,
		adding:    false,
		apiInput:  apiInput,
		fileInput: fileInput,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.adding {
		// Handle input for adding a new endpoint
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				if m.apiInput.Focused() {
					m.apiInput.Blur()
					m.fileInput.Focus()
				} else if m.fileInput.Focused() {
					// Create new endpoint and add it
					newEndpoint := model.Endpoint{
						API:          m.apiInput.Value(),
						ResponseFile: m.fileInput.Value(),
					}
					m.endpoints.AddEndpoint(newEndpoint)
					m.list.InsertItem(len(m.list.Items()), listItem(newEndpoint))

					// Reset inputs and exit adding mode
					m.apiInput.SetValue("")
					m.fileInput.SetValue("")
					m.fileInput.Blur()
					m.adding = false
				}
			case tea.KeyEsc:
				// Cancel adding
				m.apiInput.SetValue("")
				m.fileInput.SetValue("")
				m.apiInput.Blur()
				m.fileInput.Blur()
				m.adding = false
			}
		}

		// Update text inputs
		if m.apiInput.Focused() {
			m.apiInput, cmd = m.apiInput.Update(msg)
		} else if m.fileInput.Focused() {
			m.fileInput, cmd = m.fileInput.Update(msg)
		}
		return m, cmd
	}

	// Normal mode
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.adding = true
			m.apiInput.Focus()
			return m, textinput.Blink
		case "d":
			index := m.list.Index()
			if index >= 0 && index < len(m.list.Items()) {
				m.list.RemoveItem(index)
				m.endpoints.Endpoints = removeIndex(m.endpoints.Endpoints, index)
			}
			return m, nil
		case "q":
			// Save endpoints to file in style <api>:<responseFile>\n
			// Open the file for writing
			file, err := os.Create("mockzz-endpoints.txt")
			if err != nil {
				fmt.Println("Error creating file:", err)
				return m, tea.Quit
			}
			defer file.Close()

			// Iterate over the list items
			for _, listItem := range m.list.Items() {
				endpointItem, ok := listItem.(item)
				if ok {
					line := fmt.Sprintf("%s:%s\n", endpointItem.api, endpointItem.responseFile)
					_, err := file.WriteString(line)
					if err != nil {
						fmt.Println("Error writing to file:", err)
						break
					}
				}
			}
			return m, tea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.adding {
		// Display text inputs for adding a new endpoint
		return fmt.Sprintf(
			"Enter API endpoint:\n%s\n\nEnter response file path:\n%s\n\nPress Enter to continue, Esc to cancel.\n",
			m.apiInput.View(),
			m.fileInput.View(),
		)
	}
	helpText := "\nPress 'a' to add an endpoint, 'd' to delete, 'q' to save & quit.\n"
	return m.list.View() + helpText
}

func removeIndex(s []model.Endpoint, index int) []model.Endpoint {
	return append(s[:index], s[index+1:]...)
}

// Helper to convert model.Endpoint to list.Item
func listItem(e model.Endpoint) list.Item {
	return item{api: e.API, responseFile: e.ResponseFile}
}

type item struct {
	api          string
	responseFile string
}

func (i item) Title() string       { return i.api }
func (i item) Description() string { return i.responseFile }
func (i item) FilterValue() string { return i.api }
