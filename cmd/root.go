package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/mitchellh/go-wordwrap"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/wwsean08/gh-gh-status/status"
	"golang.org/x/term"
)

// stripAnsiCodes removes ANSI color codes to get actual string length
func stripAnsiCodes(s string) string {
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	return ansiRegex.ReplaceAllString(s, "")
}

// padLineToWidth pads a line (which may contain ANSI codes) to the specified width
func padLineToWidth(line string, width int) string {
	actualLength := len(stripAnsiCodes(line))
	if actualLength >= width {
		return line
	}
	return line + strings.Repeat(" ", width-actualLength)
}

// eventLoopParams contains all the parameters needed for the event loop
type eventLoopParams struct {
	client       *status.Client
	area         *pterm.AreaPrinter
	watch        bool
	sigChan      chan os.Signal
	pollChan     chan bool
	resizeChan   chan bool
	ticker       *time.Ticker
	done         chan bool
	currentState *eventLoopState
}

// eventLoopState holds the mutable state for the event loop
type eventLoopState struct {
	currentSummary *status.SystemStatus
	outputError    bool
	errMsg         string
	lastUpdate     time.Time
}

// runEventLoop executes the main event loop for handling terminal resize, polling, and rendering
func runEventLoop(params eventLoopParams) {
	for {
		select {
		case <-params.sigChan:
			// Terminal was resized, trigger immediate re-render with current data
			select {
			case params.resizeChan <- true:
			default:
				// Channel full, skip this resize event
			}
		case <-params.ticker.C:
			if params.watch {
				// Time to poll for new data
				params.pollChan <- true
			}
		case <-params.pollChan:
			// Poll for new data
			summary, err := params.client.Poll()
			if err != nil {
				params.currentState.errMsg = fmt.Sprintf("Error retrieving current GitHub status, if this is in watch mode, it will try again in 1 minute.\nError Message: %s", err.Error())
				params.currentState.outputError = true
			} else {
				params.currentState.errMsg = ""
				params.currentState.outputError = false
			}
			if summary != nil {
				params.currentState.currentSummary = summary
				params.currentState.lastUpdate = time.Now()
			}

			// Render UI with current data
			output := renderUI(params.currentState.currentSummary, params.currentState.outputError, params.currentState.errMsg, params.currentState.lastUpdate)
			params.area.Update(output)

			if !params.watch {
				params.done <- true
				return
			}
		case <-params.resizeChan:
			// Re-render with existing data (no API poll needed)
			output := renderUI(params.currentState.currentSummary, params.currentState.outputError, params.currentState.errMsg, params.currentState.lastUpdate)
			params.area.Update(output)
		}
	}
}

// renderUI generates the UI output based on current data and terminal dimensions
func renderUI(summary *status.SystemStatus, outputError bool, errMsg string, lastUpdate time.Time) string {
	// Get terminal dimensions
	termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// Fallback to default values if terminal size can't be determined
		termWidth = 80
		termHeight = 24
	}

	// Calculate available width for content (accounting for box borders and padding)
	contentWidth := termWidth - 6
	if contentWidth < 40 {
		contentWidth = 40 // Minimum width
	}

	updateTime := pterm.DefaultBasicText.Sprintf("Last Updated %s \n", lastUpdate.Format("3:04 PM"))

	var outputComponentsBox string
	var outputIncidentsBox string
	var incidentURL string
	var outputIncidents bool

	if summary != nil {
		components := summary.Components
		componentSB := strings.Builder{}

		// Build component status list with proper width
		for _, component := range components {
			if component.ID == IGNORE_GHSTATUS_COMPONENTID {
				continue
			}
			statusText := ""
			if component.Status == status.COMPONENT_OPERATIONAL {
				statusText = pterm.Green(pterm.Sprintf("%s - Operational", component.Component))
			} else if component.Status == status.COMPONENT_DEGREDADED_PERFORMANCE {
				statusText = pterm.LightYellow(pterm.Sprintf("%s - Degraded Performance", component.Component))
			} else if component.Status == status.COMPONENT_PARTIAL_OUTAGE {
				statusText = pterm.Yellow(pterm.Sprintf("%s - Partial Outage", component.Component))
			} else if component.Status == status.COMPONENT_MAJOR_OUTAGE {
				statusText = pterm.Red(pterm.Sprintf("%s - Major Outage", component.Component))
			} else {
				statusText = pterm.Sprintf("%s - %s", component.Component, component.Status)
			}
			// Pad line to ensure consistent width across all lines
			paddedLine := padLineToWidth(statusText, contentWidth)
			componentSB.WriteString(paddedLine + "\n")
		}

		incidentsSB := strings.Builder{}
		incidents := summary.Incidents
		if len(incidents) > 0 {
			incidentURL = fmt.Sprintf("https://www.githubstatus.com/incidents/%s", incidents[0].ID)
			for _, incident := range incidents[0].IncidentUpdates {
				wrappedText := wordwrap.WrapString(fmt.Sprintf("Updated %s - %s", incident.Timestamp.Local().Format("2006-01-02 3:04 PM"), incident.Update), uint(contentWidth))
				// Pad each wrapped line
				lines := strings.Split(wrappedText, "\n")
				for _, line := range lines {
					if line != "" {
						paddedLine := padLineToWidth(line, contentWidth)
						incidentsSB.WriteString(paddedLine + "\n")
					}
				}
			}
			outputIncidents = true
		}

		// Create boxes - content is now padded to terminal width
		outputComponentsBox = pterm.DefaultBox.WithTitle("System Status").WithTitleTopCenter().Sprint(componentSB.String())
		outputIncidentsBox = pterm.DefaultBox.WithTitle("Incident Updates").WithTitleTopCenter().Sprint(incidentsSB.String())
	}

	// Build output to fill terminal height
	var output strings.Builder
	output.WriteString(updateTime)

	if outputError {
		output.WriteString(errMsg)
	} else if outputIncidents {
		output.WriteString(outputComponentsBox)
		output.WriteString("\n")
		output.WriteString(incidentURL)
		output.WriteString("\n")
		output.WriteString(outputIncidentsBox)
	} else {
		output.WriteString(outputComponentsBox)
	}

	// Add padding to fill remaining terminal height
	currentLines := strings.Count(output.String(), "\n") + 1
	if currentLines < termHeight {
		padding := strings.Repeat("\n", termHeight-currentLines)
		output.WriteString(padding)
	}

	return output.String()
}

var rootCmd = &cobra.Command{
	Use:   "gh-status",
	Short: "Check the status of github.com",
	Long: `A simple command to get the current status of github.com according th githubstatus.com 
with the ability to poll it every minute to keep an eye on ongoing incidents.

To upgrade the extension run the following command:
gh extension upgrade gh-gh-status
`,
	Run: func(cmd *cobra.Command, args []string) {
		watch, err := cmd.Flags().GetBool("watch")
		if err != nil {
			log.Fatal(err)
		}
		area, _ := pterm.DefaultArea.WithFullscreen(true).Start()
		client := status.NewClient()

		// Initialize state
		state := &eventLoopState{
			currentSummary: nil,
			outputError:    false,
			errMsg:         "",
			lastUpdate:     time.Now(),
		}

		// Set up signal handler for terminal resize
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGWINCH)

		// Channels for coordinating updates
		pollChan := make(chan bool, 1)   // Channel for data polling
		resizeChan := make(chan bool, 1) // Channel for resize events
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		// Initial data poll
		pollChan <- true

		// Main event loop
		done := make(chan bool)
		params := eventLoopParams{
			client:       client,
			area:         area,
			watch:        watch,
			sigChan:      sigChan,
			pollChan:     pollChan,
			resizeChan:   resizeChan,
			ticker:       ticker,
			done:         done,
			currentState: state,
		}
		go runEventLoop(params)

		// Wait for completion
		<-done
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
