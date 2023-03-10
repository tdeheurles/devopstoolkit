from typing import Dict


class Usage:
    def __init__(self, definition: str, family: str, commands: Dict[str, Dict]):
        self.definition = definition
        self.family = family
        self.commands = commands
