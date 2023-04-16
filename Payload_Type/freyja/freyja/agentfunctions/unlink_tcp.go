package agentfunctions

import (
	"fmt"
	agentstructs "github.com/MythicMeta/MythicContainer/agent_structs"
)

func init() {
	agentstructs.AllPayloadData.Get("freyja").AddCommand(agentstructs.Command{
		Name:                "unlink_tcp",
		Description:         "Unlink a freyja_tcp connection.",
		HelpString:          "unlink_tcp",
		Version:             1,
		MitreAttackMappings: []string{},
		SupportedUIFeatures: []string{},
		Author:              "@its_a_feature_",
		CommandAttributes: agentstructs.CommandAttribute{
			CommandIsSuggested: false,
			CommandIsBuiltin:   false,
      SupportedOS: []string{agentstructs.SUPPORTED_OS_LINUX, agentstructs.SUPPORTED_OS_MACOS, agentstructs.SUPPORTED_OS_WINDOWS},
		},
		CommandParameters: []agentstructs.CommandParameter{
			{
				Name:          "connection",
				Description:   "Connection info for unlinking",
				ParameterType: agentstructs.COMMAND_PARAMETER_TYPE_LINK_INFO,
			},
		},
		TaskFunctionParseArgString: func(args *agentstructs.PTTaskMessageArgsData, input string) error {
			return nil
		},
		TaskFunctionParseArgDictionary: func(args *agentstructs.PTTaskMessageArgsData, input map[string]interface{}) error {
			return args.LoadArgsFromDictionary(input)
		},
		TaskFunctionCreateTasking: func(taskData *agentstructs.PTTaskMessageAllData) agentstructs.PTTaskCreateTaskingMessageResponse {
			response := agentstructs.PTTaskCreateTaskingMessageResponse{
				Success: true,
				TaskID:  taskData.Task.ID,
			}
			if connectionInfo, err := taskData.Args.GetDictionaryArg("connection"); err != nil {
				response.Success = false
				response.Error = err.Error()
			} else if callbackUUID, ok := connectionInfo["callback_uuid"]; !ok {
				response.Success = false
				response.Error = "Failed to find callback UUID in connection information"
			} else {
				taskData.Args.RemoveArg("connection")
				taskData.Args.AddArg(agentstructs.CommandParameter{
					Name:          "connection",
					DefaultValue:  callbackUUID,
					ParameterType: agentstructs.COMMAND_PARAMETER_TYPE_STRING,
				})
				displayString := fmt.Sprintf("from %s", callbackUUID)
				response.DisplayParams = &displayString
			}
			return response
		},
	})
}
