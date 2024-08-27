module github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code

go 1.19

require (
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	golang.org/x/sys v0.18.0
)

require golang.org/x/net v0.23.0 // indirect

//replace github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code => ./
