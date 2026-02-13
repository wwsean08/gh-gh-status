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
	"unsafe"

	"github.com/mitchellh/go-wordwrap"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/wwsean08/gh-gh-status/status"
	"golang.org/x/term"
)

// setNonCanonicalMode sets the terminal to non-canonical mode for reading
// single characters without affecting output processing
func setNonCanonicalMode(fd int) (*term.State, error) {
	oldState, err := term.GetState(fd)
	if err != nil {
		return nil, err
	}

	// Get the raw terminal attributes
	var termios syscall.Termios
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCGETA, uintptr(unsafe.Pointer(&termios))); errno != 0 {
		return oldState, errno
	}

	// Modify only input flags - leave output flags untouched
	// Disable canonical mode (ICANON) and echo (ECHO)
	termios.Lflag &^= syscall.ICANON | syscall.ECHO | syscall.ECHOE | syscall.ECHOK | syscall.ECHONL
	// Set minimum characters to read
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0

	// Apply the modified settings
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), syscall.TIOCSETA, uintptr(unsafe.Pointer(&termios))); errno != 0 {
		return oldState, errno
	}

	return oldState, nil
}

// handleKeyboardInput listens for keyboard input and sends signals to appropriate channels
// Note: Terminal must already be in non-canonical mode before calling this function
func handleKeyboardInput(refreshChan chan bool, done chan bool) {
	buf := make([]byte, 1)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil || n == 0 {
			continue
		}

		switch buf[0] {
		case 'r', 'R':
			// Manual refresh requested
			select {
			case refreshChan <- true:
			default:
				// Channel full, skip
			}
		case 'q', 'Q', 3: // 3 is Ctrl+C
			// Quit requested
			select {
			case done <- true:
			default:
			}
			return
		}
	}
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// stripAnsiCodes removes ANSI color codes to get actual string length
func stripAnsiCodes(s string) string {
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
	refreshChan  chan bool
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
		case <-params.refreshChan:
			// Manual refresh requested (in watch mode only)
			if params.watch {
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
				// Update last check time even if there's no new data (304 response)
				params.currentState.lastUpdate = time.Now()
			}
			if summary != nil {
				params.currentState.currentSummary = summary
			}

			// Render UI with current data
			output := renderUI(params.currentState.currentSummary, params.currentState.outputError, params.currentState.errMsg, params.currentState.lastUpdate, params.watch)
			params.area.Update(output)

			if !params.watch {
				params.done <- true
				return
			}
		case <-params.resizeChan:
			// Clear the area to prevent artifacts from previous render
			params.area.Clear()
			// Re-render with existing data (no API poll needed)
			output := renderUI(params.currentState.currentSummary, params.currentState.outputError, params.currentState.errMsg, params.currentState.lastUpdate, params.watch)
			params.area.Update(output)
		}
	}
}

// renderUI generates the UI output based on current data and terminal dimensions
func renderUI(summary *status.SystemStatus, outputError bool, errMsg string, lastUpdate time.Time, watch bool) string {
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
	linesNeeded := termHeight - currentLines

	if watch {
		// In watch mode, add help text at bottom with blank line separator
		// Calculate how much space we have: termHeight - currentLines - 1 (for help text)
		linesNeeded -= 1
		if linesNeeded > 0 {
			// Add padding + blank line separator
			output.WriteString(strings.Repeat("\n", linesNeeded))
		}
		// Add help text at the bottom (no trailing newline)
		output.WriteString(pterm.Gray("Press 'r' to refresh, 'q' or Ctrl+C to quit"))
	} else {
		// Normal mode: fill to terminal height
		if currentLines < termHeight {
			padding := strings.Repeat("\n", linesNeeded)
			output.WriteString(padding)
		}
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

		// Set up signal handler for terminal resize and interrupt
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGWINCH)

		intChan := make(chan os.Signal, 1)
		signal.Notify(intChan, os.Interrupt, syscall.SIGTERM)

		// Channels for coordinating updates
		pollChan := make(chan bool, 1)    // Channel for data polling
		resizeChan := make(chan bool, 1)  // Channel for resize events
		refreshChan := make(chan bool, 1) // Channel for manual refresh
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		// Initial data poll
		pollChan <- true

		// Main event loop
		done := make(chan bool, 1)

		// Set up terminal for keyboard input in watch mode
		// Use non-canonical mode instead of full raw mode to preserve output processing
		var oldTermState *term.State
		if watch {
			oldTermState, err = setNonCanonicalMode(int(os.Stdin.Fd()))
			if err == nil {
				defer term.Restore(int(os.Stdin.Fd()), oldTermState)
				// Start keyboard input handler
				go handleKeyboardInput(refreshChan, done)
			}
		}

		params := eventLoopParams{
			client:       client,
			area:         area,
			watch:        watch,
			sigChan:      sigChan,
			pollChan:     pollChan,
			resizeChan:   resizeChan,
			refreshChan:  refreshChan,
			ticker:       ticker,
			done:         done,
			currentState: state,
		}
		go runEventLoop(params)

		// Handle interrupt signal
		go func() {
			<-intChan
			done <- true
		}()

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
