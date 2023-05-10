package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/wwsean08/gh-status/status"
	"log"
)

var rootCmd = &cobra.Command{
	Use:   "gh-status",
	Short: "Check the status of github.com",
	Long: `A simple command to get the current status of github.com according th githubstatus.com 
with the ability to poll it every minute to keep an eye on ongoing incidents.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := status.NewClient()
		status, err := client.Poll()
		if err != nil {
			log.Fatal(err)
		}
		if status != nil {
			s, _ := json.Marshal(status)
			fmt.Printf("success: \n%s\n", s)
		} else {
			log.Fatal("status is nil")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().BoolP("watch", "w", false, "Check for a status update every minute")
}
