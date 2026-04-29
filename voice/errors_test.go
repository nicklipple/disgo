package voice

import (
	"errors"
	"testing"

	"github.com/disgoorg/snowflake/v2"
)

func TestVoiceError_Error(t *testing.T) {
	guildID := snowflake.ID(123456789)

	tests := []struct {
		name       string
		voiceErr   VoiceError
		want       string
		wantUnwrap string
		wantNilErr bool
	}{
		{
			name: "with underlying error",
			voiceErr: VoiceError{
				GuildID:     guildID,
				Code:        4004,
				Description: "Authentication failed",
				Explanation: "The token you sent in your identify payload is incorrect.",
				Resumeable:  false,
				Err:         errors.New("websocket close 4004"),
			},
			want:       "voice error (guild: 123456789, code: 4004): Authentication failed - The token you sent in your identify payload is incorrect.: websocket close 4004",
			wantUnwrap: "websocket close 4004",
			wantNilErr: false,
		},
		{
			name: "without underlying error",
			voiceErr: VoiceError{
				GuildID:     guildID,
				Code:        4004,
				Description: "Authentication failed",
				Explanation: "The token you sent in your identify payload is incorrect.",
				Resumeable:  false,
				Err:         nil,
			},
			want:       "voice error (guild: 123456789, code: 4004): Authentication failed - The token you sent in your identify payload is incorrect.",
			wantUnwrap: "",
			wantNilErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.voiceErr.Error(); got != tt.want {
				t.Errorf("VoiceError.Error() = %v, want %v", got, tt.want)
			}

			unwrapped := tt.voiceErr.Unwrap()
			if tt.wantNilErr {
				if unwrapped != nil {
					t.Errorf("VoiceError.Unwrap() = %v, want nil", unwrapped)
				}
			} else {
				if unwrapped == nil {
					t.Errorf("VoiceError.Unwrap() = nil, want %v", tt.wantUnwrap)
				} else if unwrapped.Error() != tt.wantUnwrap {
					t.Errorf("VoiceError.Unwrap() = %v, want %v", unwrapped.Error(), tt.wantUnwrap)
				}
			}
		})
	}
}
