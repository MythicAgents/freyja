package agentfunctions

import (
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"
	"github.com/MythicMeta/MythicContainer/mythicrpc"
)

var powershell_executor = agentstructs.Command{
	Name:                      "powershell_executor",
  Author:                    "@antman1p",
	Description:               "Execute a shell command using 'powershell -nologo -noprofile'",
	MitreAttackMappings:       []string{"T1059"},
	CommandAttributes: agentstructs.CommandAttribute{
		CommandIsSuggested: true,
		CommandIsBuiltin: false,
		FilterCommandAvailabilityByAgentBuildParameters: true,
		SupportedOS: []string{agentstructs.SUPPORTED_OS_WINDOWS},
	},
	TaskFunctionCreateTasking: powershell_executorCreateTasking,
	TaskFunctionParseArgDictionary: func(args *agentstructs.PTTaskMessageArgsData, input map[string]interface{}) error {
		return args.LoadArgsFromDictionary(input)
	},
	Version: 1,
}

func init() {
	agentstructs.AllPayloadData.Get("freyja").AddCommand(powershell_executor)
}

func powershell_executorCreateTasking(taskData *agentstructs.PTTaskMessageAllData) agentstructs.PTTaskCreateTaskingMessageResponse {
	response := agentstructs.PTTaskCreateTaskingMessageResponse{
		Success: true,
		TaskID:  taskData.Task.ID,
	}
	if _, err := mythicrpc.SendMythicRPCArtifactCreate(mythicrpc.MythicRPCArtifactCreateMessage{
		BaseArtifactType: "ProcessCreate",
		ArtifactMessage:  "powershell -nologo -noprofile " + taskData.Args.GetCommandLine(),
		TaskID:           taskData.Task.ID,
	}); err != nil {
		logging.LogError(err, "Failed to send mythicrpc artifact create")
	}
	return response
}
