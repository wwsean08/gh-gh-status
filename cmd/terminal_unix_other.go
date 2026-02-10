//go:build linux

package cmd

import (
	"golang.org/x/sys/unix"
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
	termios, err := unix.IoctlGetTermios(fd, unix.TCGETS)
	if err != nil {
		return oldState, err
	}

	// Modify only input flags - leave output flags untouched
	// Disable canonical mode (ICANON) and echo (ECHO)
	termios.Lflag &^= unix.ICANON | unix.ECHO | unix.ECHOE | unix.ECHOK | unix.ECHONL
	// Set minimum characters to read
	termios.Cc[unix.VMIN] = 1
	termios.Cc[unix.VTIME] = 0

	// Apply the modified settings
	if err := unix.IoctlSetTermios(fd, unix.TCSETS, termios); err != nil {
		return oldState, err
	}

	return oldState, nil
}
