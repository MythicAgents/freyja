package agentfunctions

import (
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/mythicrpc"
)

var bash_executor = agentstructs.Command{
	Name:                      "bash_executor",
  Author:                    "@antman1p",
	Description:               "Execute a shell command using 'bash -c'",
	MitreAttackMappings:       []string{"T1059"},
  CommandAttributes: agentstructs.CommandAttribute{
      CommandIsSuggested: true,
      SupportedOS: []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_MACOS},
  },
	TaskFunctionCreateTasking: bash_executorCreateTasking,
	TaskFunctionParseArgDictionary: func(args *agentstructs.PTTaskMessageArgsData, input map[string]interface{}) error {
		return args.LoadArgsFromDictionary(input)
	},
	Version: 1,
}

func init() {
	agentstructs.AllPayloadData.Get("freyja").AddCommand(bash_executor)
}

func bash_executorCreateTasking(taskData *agentstructs.PTTaskMessageAllData) agentstructs.PTTaskCreateTaskingMessageResponse {
	response := agentstructs.PTTaskCreateTaskingMessageResponse{
		Success: true,
		TaskID:  taskData.Task.ID,
	}
	if _, err := mythicrpc.SendMythicRPCArtifactCreate(mythicrpc.MythicRPCArtifactCreateMessage{
		BaseArtifactType: "ProcessCreate",
		ArtifactMessage:  "/bin/bash -c " + taskData.Args.GetCommandLine(),
		TaskID:           taskData.Task.ID,
	}); err != nil {
		logging.LogError(err, "Failed to send mythicrpc artifact create")
	}
	return response
}
