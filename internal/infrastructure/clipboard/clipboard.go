// Package clipboard provides clipboard operations.
package clipboard

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/DementevVV/commitsum/internal/domain/repository"
)

// Service provides clipboard operations.
type Service struct{}

// Ensure Service implements ClipboardRepository.
var _ repository.ClipboardRepository = (*Service)(nil)

// New creates a new clipboard service.
func New() *Service {
	return &Service{}
}

// Copy copies text to the system clipboard.
func (s *Service) Copy(text string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		// Try xclip first, then xsel, then wl-copy (Wayland).
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("xsel"); err == nil {
			cmd = exec.Command("xsel", "--clipboard", "--input")
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			// Fallback to xclip anyway, let it fail with proper error.
			cmd = exec.Command("xclip", "-selection", "clipboard")
		}
	case "windows":
		cmd = exec.Command("cmd", "/c", "clip")
	default:
		cmd = exec.Command("pbcopy") // Default to macOS.
	}

	cmd.Stdin = strings.NewReader(text)
	return cmd.Run()
}

// IsAvailable checks if clipboard is available on the system.
func (s *Service) IsAvailable() bool {
	cmd := s.getCommand()
	_, err := exec.LookPath(cmd)
	return err == nil
}

func (s *Service) getCommand() string {
	switch runtime.GOOS {
	case "darwin":
		return "pbcopy"
	case "linux":
		if _, err := exec.LookPath("xclip"); err == nil {
			return "xclip"
		} else if _, err := exec.LookPath("xsel"); err == nil {
			return "xsel"
		} else if _, err := exec.LookPath("wl-copy"); err == nil {
			return "wl-copy"
		}
		return "xclip"
	case "windows":
		return "clip"
	default:
		return "pbcopy"
	}
}
