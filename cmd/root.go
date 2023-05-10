package cmd

import (
	"fmt"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/wwsean08/gh-status/status"
	"log"
	"strings"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "gh-status",
	Short: "Check the status of github.com",
	Long: `A simple command to get the current status of github.com according th githubstatus.com 
with the ability to poll it every minute to keep an eye on ongoing incidents.`,
	Run: func(cmd *cobra.Command, args []string) {
		area, _ := pterm.DefaultArea.WithFullscreen().Start()
		client := status.NewClient()
		outputComponentsBox := ""
		outputIncidentsBox := ""
		outputIncidents := false
		for {
			updateTime := pterm.DefaultBasicText.Sprintf(fmt.Sprintf("Last Updated %s\n", time.Now().Format(time.TimeOnly)))
			status, err := client.Poll()
			if err != nil {
				log.Fatal(err)
			}
			if status != nil {
				components := status.Components
				componentSB := strings.Builder{}
				for _, component := range components {
					if component.ID == IGNORE_GHSTATUS_COMPONENTID {
						continue
					}
					if component.Status == "operational" {
						componentSB.WriteString(pterm.Green(pterm.Sprintf("%s - Operational\n", component.Component)))
					} else if component.Status == "partial_outage" {
						componentSB.WriteString(pterm.Yellow(pterm.Sprintf("%s - Degraded\n", component.Component)))
					}
				}

				incidentsSB := strings.Builder{}
				incidents := status.Incidents
				if len(incidents) > 0 {
					for _, incident := range incidents[0].IncidentUpdates {
						incidentsSB.WriteString(fmt.Sprintf("Updated At %s - %s\n", incident.Timestamp.Local().Format(time.DateTime), incident.Update))
					}
					outputIncidents = true
				} else {
					outputIncidents = false
				}
				outputComponentsBox = pterm.DefaultBox.WithTitle("System Status").WithTitleTopCenter().Sprint(componentSB.String())
				outputIncidentsBox = pterm.DefaultBox.WithTitle("Incident Updates").WithTitleTopCenter().Sprint(incidentsSB.String())
			}
			if outputIncidents {
				area.Update(updateTime, outputComponentsBox, "\n\n", outputIncidentsBox)
			} else {
				area.Update(updateTime, outputComponentsBox)

			}
			time.Sleep(time.Minute)
		}
		area.Stop()
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
