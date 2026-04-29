package voice

import (
	"fmt"

	"github.com/disgoorg/snowflake/v2"
)

// VoiceError represents a voice gateway error with rich context information.
type VoiceError struct {
	// GuildID is the ID of the guild where the error occurred.
	GuildID snowflake.ID

	// Code is the voice gateway close event code (e.g., 4004).
	Code int

	// Description is the short description of the error (e.g., "Authentication failed").
	Description string

	// Explanation provides more details about why the error occurred.
	Explanation string

	// Resumeable indicates whether the connection can attempt to resume.
	Resumeable bool

	// Err is the underlying error that caused the close (may be nil for some close codes).
	Err error
}

// Error implements the error interface.
func (e VoiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("voice error (guild: %s, code: %d): %s - %s: %v",
			e.GuildID, e.Code, e.Description, e.Explanation, e.Err)
	}
	return fmt.Sprintf("voice error (guild: %s, code: %d): %s - %s",
		e.GuildID, e.Code, e.Description, e.Explanation)
}

// Unwrap returns the underlying error.
func (e VoiceError) Unwrap() error {
	return e.Err
}
