from mythic_payloadtype_container.MythicCommandBase import *
from mythic_payloadtype_container.MythicRPC import *


class ShellArguments(TaskArguments):
    def __init__(self, command_line, **kwargs):
        super().__init__(command_line, **kwargs)
        self.args = [
            CommandParameter(name="command", display_name="Command", type=ParameterType.String, description="Command to run"),
        ]

    async def parse_arguments(self):
        if len(self.command_line) == 0:
            raise ValueError("Must supply a command to run")
        self.add_arg("command", self.command_line)

    async def parse_dictionary(self, dictionary_arguments):
        self.load_args_from_dictionary(dictionary_arguments)


class ShellCommand(CommandBase):
    cmd = "zsh"
    needs_admin = False
    help_cmd = "zsh {command}"
    description = "Execute a shell command using 'zsh -c'"
    version = 1
    author = "@antman1p"
    attackmapping = ["T1059", "T1059.004"]
    argument_class = ShellArguments
    attributes = CommandAttributes(
        suggested_command=True
    )

    async def create_tasking(self, task: MythicTask) -> MythicTask:
        resp = await MythicRPC().execute("create_artifact", task_id=task.id,
            artifact="/bin/zsh -c {}".format(task.args.get_arg("command")),
            artifact_type="Process Create",
        )
        resp = await MythicRPC().execute("create_artifact", task_id=task.id,
            artifact="{}".format(task.args.get_arg("command")),
            artifact_type="Process Create",
        )
        task.display_params = task.args.get_arg("command")
        return task

    async def process_response(self, response: AgentResponse):
        pass
