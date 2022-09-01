from mythic_payloadtype_container.MythicCommandBase import *
import json
from mythic_payloadtype_container.MythicRPC import *


class ZshArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = []

    async def parse_arguments(self):
        pass

class ZshCommand(CommandBase):
    cmd = "zsh_executor"
    needs_admin = False
    help_cmd = "zsh {command}"
    description = "Execute a shell command using 'zsh -c'"
    version = 1
    author = "@antman1p"
    attackmapping = ["T1059", "T1059.004"]
    argument_class = ZshArguments
    attributes = CommandAttributes(
        suggested_command=True
    )

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        resp = await MythicRPC().execute("create_artifact", task_id=task.id,
            artifact="/bin/zsh -c {}".format(task.args.command_line),
            artifact_type="Process Create",
        )
        resp = await MythicRPC().execute("create_artifact", task_id=task.id,
            artifact="{}".format(task.args.command_line),
            artifact_type="Process Create",
        )
        #task.display_params = task.args.get_arg("command")
        return task

    async def process_response(self, response: AgentResponse):
        pass
