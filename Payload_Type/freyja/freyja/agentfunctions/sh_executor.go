package agentfunctions

import (
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/mythicrpc"
)

var sh_executor = agentstructs.Command{
	Name:                      "sh_executor",
  Author:                    "@antman1p",
	Description:               "Execute a shell command using 'sh -c'",
	MitreAttackMappings:       []string{"T1059"},
  CommandAttributes: agentstructs.CommandAttribute{
      CommandIsSuggested: true,
      SupportedOS: []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_MACOS},
  },
	TaskFunctionCreateTasking: sh_executorCreateTasking,
	TaskFunctionParseArgDictionary: func(args *agentstructs.PTTaskMessageArgsData, input map[string]interface{}) error {
		return args.LoadArgsFromDictionary(input)
	},
	Version: 1,
}

func init() {
	agentstructs.AllPayloadData.Get("freyja").AddCommand(sh_executor)
}

func sh_executorCreateTasking(taskData *agentstructs.PTTaskMessageAllData) agentstructs.PTTaskCreateTaskingMessageResponse {
	response := agentstructs.PTTaskCreateTaskingMessageResponse{
		Success: true,
		TaskID:  taskData.Task.ID,
	}
	if _, err := mythicrpc.SendMythicRPCArtifactCreate(mythicrpc.MythicRPCArtifactCreateMessage{
		BaseArtifactType: "ProcessCreate",
		ArtifactMessage:  "/bin/sh -c " + taskData.Args.GetCommandLine(),
		TaskID:           taskData.Task.ID,
	}); err != nil {
		logging.LogError(err, "Failed to send mythicrpc artifact create")
	}
	return response
}
