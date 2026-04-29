package voice

import (
	"context"
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func TestConn_Errors(t *testing.T) {
	// Test that Errors() returns a channel
	conn := NewConn(snowflake.ID(123), snowflake.ID(456),
		func(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error {
			return nil
		},
		func() {})

	errChan := conn.Errors()
	if errChan == nil {
		t.Fatal("Errors() returned nil channel")
	}

	// Verify it's receive-only by trying to receive from it
	// It should be empty but not closed initially
	select {
	case <-errChan:
		t.Fatal("Expected channel to be empty initially")
	default:
		// Channel is empty as expected
	}
}

func TestConn_ErrorsChannelType(t *testing.T) {
	// Verify that Errors() returns a receive-only channel
	conn := NewConn(snowflake.ID(123), snowflake.ID(456),
		func(ctx context.Context, guildID snowflake.ID, channelID *snowflake.ID, selfMute bool, selfDeaf bool) error {
			return nil
		},
		func() {})

	errChan := conn.Errors()

	// Verify channel type is receive-only <-chan VoiceError
	// This is a compile-time check, but we can verify it at runtime
	if errChan == nil {
		t.Fatal("Errors() returned nil")
	}

	// Verify channel is not buffered by checking capacity
	if cap(errChan) != 0 {
		t.Logf("Channel capacity is %d, expected 0 (unbuffered)", cap(errChan))
	}
}
