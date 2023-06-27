package cmd

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/wwsean08/gh-gh-status/status"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "gh-status",
	Short: "Check the status of github.com",
	Long: `A simple command to get the current status of github.com according th githubstatus.com 
with the ability to poll it every minute to keep an eye on ongoing incidents.

To upgrade the extension run the following command:
gh extension upgrade gh-gh-status
`,
	Run: func(cmd *cobra.Command, args []string) {
		m := model{
			components: list.New(nil, NewItemDelegator(), 0, 0),
		}
		p := tea.NewProgram(
			m,
			tea.WithAltScreen(), // use the full size of the terminal in its "alternate screen buffer"
		)

		client := status.NewClient()
		go func(program tea.Program) {
			for {
				status, err := client.Poll()
				timestamp := time.Now()
				msg := statusMsg{
					status,
					err,
					&timestamp,
				}
				p.Send(msg)
				time.Sleep(time.Minute)
			}
		}(*p)

		if _, err := p.Run(); err != nil {
			fmt.Println("could not run program:", err)
			os.Exit(1)
		}

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
