# voice

Voice provides a package to connect and send/receive voice to/from discord servers.
For Discords Docs on voice see [here](https://discord.com/developers/docs/topics/voice-connections).

DAVE(E2EE) library options:
 * https://github.com/disgoorg/godave CGO binding for https://github.com/discord/libdave
 * https://github.com/thomas-vilte/dave-go Pure Go implementation of DAVE(E2EE) (experimental)

## GoDave

### Installation

```bash
go get github.com/disgoorg/godave/golibdave
```

### Logging
Libdave uses a global logger which is set it `slog.LevelError` by default. You can change this by calling:

```go
libdave.SetDefaultLogLoggerLevel(slog.LevelInfo)
```

or set your own logger:

```go
libdave.SetDefaultLogLogger(yourLogger)
```

## Dave-Go

### Installation

```bash
go get github.com/thomas-vilte/dave-go
```

## Usage

To send audio you need to create a voice connection. When using the `bot.Client` package you can use `client.VoiceManager().CreateConn(guildID)`
```go
const (
    guildID = 12345
    channelID = 12345
)

client, err := disgo.New(token,
	bot.WithGatewayConfigOpts(gateway.WithIntents(gateway.IntentGuildVoiceStates)),
	bot.WithVoiceManagerConfigOpts(
		// for GoDave use this
		voice.WithDaveSessionCreateFunc(golibdave.NewSession),
		// for Dave-Go use this
		voice.WithDaveSessionCreateFunc(session.NewSession),
	),
)
// handle err

conn := client.VoiceManager().CreateConn(guildID)

err := conn.Open(context.TODO(), channelID, false, false)
// handle err

// set speaking flag
err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone)

// send opus frame
conn.UDP().Write(frame)

// close connection
conn.Close()
```

## Error Handling

Voice connections expose errors through a synchronous channel. This allows callers to handle non-resumeable errors by closing and restarting connections.

```go
conn := client.VoiceManager.CreateConn(guildID)

// Subscribe to errors
go func() {
    for err := range conn.Errors() {
        voiceErr := err.(voice.VoiceError)
        
        log.Printf("Voice error: code=%d, description=%s, resumeable=%v",
            voiceErr.Code, voiceErr.Description, voiceErr.Resumeable)
        
        if !voiceErr.Resumeable {
            // Non-resumeable error - close and restart
            conn.Close(ctx)
            conn = client.VoiceManager.CreateConn(guildID)
            if err := conn.Open(ctx, channelID, false, false); err != nil {
                log.Fatal(err)
            }
        }
    }
}()
```

**Important Notes:**
- The error channel is unbuffered and synchronous
- Only one subscriber should range over `Errors()`
- The channel remains open during auto-reconnects; it is garbage collected when the Conn is garbage collected
- All voice gateway close codes are exposed (both resumeable and non-resumeable)

When using the voice package standalone you should create a voice manager. After this you can call `voice.Manager.CreateConn(guildID)`. After this you should send a `gateway.OpcodeVoiceStateUpdate` packet to the gateway.
```go
