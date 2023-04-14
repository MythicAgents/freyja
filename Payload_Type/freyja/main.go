package main

import (
	freyjafunctions "MyContainer/freyja/agentfunctions"
	poseidontcpfunctions "MyContainer/freyja_tcp/c2functions"
	"github.com/MythicMeta/MythicContainer"
)

func main() {
	// load up the agent functions directory so all the init() functions execute
	freyjafunctions.Initialize()
	freyjatcpfunctions.Initialize()
	// sync over definitions and listen
	MythicContainer.StartAndRunForever([]MythicContainer.MythicServices{
		MythicContainer.MythicServicePayload,
		MythicContainer.MythicServiceC2,
	})
}
