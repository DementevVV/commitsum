package repository

// ClipboardRepository defines the interface for clipboard operations.
type ClipboardRepository interface {
	// Copy copies text to the system clipboard.
	Copy(text string) error

	// IsAvailable checks if clipboard is available on the system.
	IsAvailable() bool
}
