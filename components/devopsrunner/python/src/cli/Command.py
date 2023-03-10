from typing import List


class Command:
    def __init__(
            self,
            name: str,
            execution,
            help: str = "undefined",
            command_example: List[str] = None,
            parsed_files: List[str] = None,
            files: List[str] = None,
            before_parse=None):
        # All Commands
        self.name = name
        self.before_parse = before_parse if before_parse is not None else lambda *args: None
        self.execution = execution
        self.help = help
        self.command_example = command_example if command_example is not None else []

        # File System Commands
        self.parsed_files = parsed_files if parsed_files is not None else []
        self.files = files if files is not None else []
        self.template_files = []

    def __str__(self) -> str:
        return f" -name:{self.name} -files: {self.files}"
