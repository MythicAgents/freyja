module github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code

go 1.21

require (
	github.com/creack/pty v1.1.21
	github.com/djherbis/atime v1.1.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/kbinani/screenshot v0.0.0-20230812210009-b87d31814237
	github.com/tmc/scp v0.0.0-20170824174625-f7b48647feef
	github.com/xorrior/keyctl v1.0.1-0.20210425144957-8746c535bf58
	golang.org/x/crypto v0.20.0
	golang.org/x/sync v0.6.0
	golang.org/x/sys v0.17.0
	howett.net/plist v1.0.1
)

require (
	github.com/MythicAgents/freyja v0.0.0-20240309224810-7454af137fd9 // indirect
	github.com/gen2brain/shm v0.1.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/lxn/win v0.0.0-20210218163916-a377121e959e // indirect
	golang.org/x/net v0.21.0 // indirect
)

replace github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code => ./
