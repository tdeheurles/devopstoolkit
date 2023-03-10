import datetime
import os
import shutil
from typing import List, Dict
from sys import exit
from cli.Command import Command
from cli.TemplateFile import TemplateFile
from cli.log import cli_print, debug_print
from cli.parameters import indentation, key_script_name, key_component, \
    key_command, keyword_command, key_command_short, print_usage_options, key, \
    get_and_assert_additional_parameters
from cli.theme import theme_main0, theme_main1, theme_error0

key_runner_file_prefix = "runner."
key_help = "help"
key_files = "files"


class Component:
    def __init__(self, name: str, definition: str, family: str, commands: Dict[str, Command], directory: str = None):
        self.internal_component = True if directory is None else False
        self.name = name
        self.directory = directory
        self.definition = definition
        self.family = family
        self.commands = commands

    def __str__(self) -> str:
        to_str = f"Component: {self.name}\n" \
                 f" -directory:{self.directory}\n" \
                 f" -definition:{self.definition}\n" \
                 f" -family:{self.family}\n" \
                 f" -commands:\n"

        for command_name, command_data in self.commands.items():
            to_str += f"  {str(command_data)}\n"

        return to_str

    def print_component_usage(self):
        print()
        cli_print(f"Component: {self.name}", theme_main0)
        cli_print(self.definition, theme_main1)

        print()
        cli_print('Usage:', theme_main0)
        cli_print(
            f'{indentation}{key_script_name}'
            f' --{key_component} {self.name}'
            f' --{key_command} {keyword_command}'
            f' [Options]'
        )

        print()
        cli_print(f'Where (--{key_command}|-{key_command_short}) {keyword_command} is one of:', theme_main0)
        max_length = 15
        for command_name, command in self.commands.items():
            if len(command_name) > max_length:
                max_length = len(command_name)
        command_format = indentation + "{0:<" + f"{max_length + 3}" + "}{1:<15}"
        command_names: List[str] = list(self.commands.keys())
        command_names.sort()
        for command_name in command_names:
            cli_print(command_format.format(command_name, self.commands[command_name].help))

        print()
        print_usage_options()

    @staticmethod
    def print_command_usage(component_name: str, command: Command):

        print()
        cli_print(f"Component: {component_name} - Command: {command.name}", theme_main0)
        cli_print(command.help, theme_main1)

        print()
        cli_print('Usage:', theme_main0)
        params_from_files = list(set([
            f"--{param} {param.upper()}"
            for template_file in command.template_files
            for param in template_file.required_params
        ]))

        command_base = f'{indentation}{key_script_name}' + \
                       f' --{key_component} {component_name}'
        cli_print(
            command_base +
            f' --{key_command} {command.name}'
            f' {" ".join(params_from_files)}'
            f' [Options]'
        )

        if len(command.command_example) != 0:
            print()
            cli_print("Examples:", theme_main0)
            for example in command.command_example:
                cli_print(f"{indentation}{command_base} {example}")

        print()
        print_usage_options()

    @staticmethod
    def file_without_runner_prefix(file_name: str) -> str:
        return file_name.replace(key_runner_file_prefix, "")

    @staticmethod
    def file_with_runner_prefix(file_name: str) -> str:
        return f"{key_runner_file_prefix}{file_name}"

    def execute_command(self, params: Dict[str, str]) -> int:
        debug_print("\n\nExecute command")

        command_name = params[key_command]
        command = self.commands[command_name]

        command.before_parse()
        ts00 = datetime.datetime.now().timestamp()
        self.parse_files(command_name=command_name, params=params)
        ts01 = datetime.datetime.now().timestamp()
        debug_print(f"parse_files: {ts01 - ts00}")

        ts10 = datetime.datetime.now().timestamp()
        exit_code = command.execution(params, self, command)
        ts11 = datetime.datetime.now().timestamp()
        debug_print(f"exec: {ts11 - ts10}")

        ts20 = datetime.datetime.now().timestamp()
        self.remove_files(command_name=command_name)
        ts21 = datetime.datetime.now().timestamp()
        debug_print(f"remove: {ts21 - ts20}")
        
        return exit_code >> 8  # os.system exit code return 256 instead of 1

    def parse_files(self, command_name: str, params: Dict[str, str]):
        debug_print("--parse_files-->")
        command = self.commands[command_name]

        # Get required parameters
        template_files = self.get_file_templates(command_name=command_name)
        command.template_files = template_files
        debug_print(f"|_ template_files: {command.template_files}")

        all_params = get_and_assert_additional_parameters(
            initial_params=params,
            usage=lambda cn=self.name, cd=command: self.print_command_usage(cn, cd),
            new_required_params=[
                {key: required_param}
                for template_file in template_files
                for required_param in template_file.required_params])

        for template_file in template_files:
            debug_print(f"--parsing file--> {template_file}")
            parsed_file = self.file_without_runner_prefix(file_name=template_file.path)
            command.parsed_files.append(parsed_file)

            with open(self.get_full_path(parsed_file), 'w') as to_stream:
                new_content = template_file.content
                for required_parameter in template_file.required_params:
                    new_content = new_content.replace(
                        f"<% {required_parameter} %>",
                        all_params[required_parameter])

                to_stream.write(new_content)
            if ".sh" in self.get_full_path(parsed_file) or ".py" in self.get_full_path(parsed_file):
                os.system(f"chmod u+x {self.get_full_path(parsed_file)}")

    def remove_files(self, command_name: str):
        for file in self.commands[command_name].parsed_files:
            self.remove_file(file)

    def remove_file(self, file: str):
        try:
            os.remove(self.get_full_path(file))
        except FileNotFoundError:
            pass

    def get_full_path(self, file_name) -> str:
        return f"{self.directory}/{file_name}"

    def get_file_templates(self, command_name) -> List[TemplateFile]:
        template_files = []
        for template_file in self.commands[command_name].files:
            with open(self.get_full_path(template_file), 'r') as from_stream:
                template_files.append(TemplateFile(
                    template_file=template_file,
                    template_content=from_stream.read()))

        return template_files

    @staticmethod
    def execute_filesystem_before_parse(from_file: str, to_file: str) -> int:
        debug_print(f"--before_parse--> copying {from_file} to {to_file}")
        shutil.copy(from_file, to_file)
        return 0

    @staticmethod
    def execute_filesystem_command(path) -> int:
        if path[-3:] == ".go":
            filename = os.path.basename(path)
            dirname = os.path.dirname(path)
            return os.system(f"cd {dirname} && go run {filename}")
        return os.system(f"chmod u+x {path} && cd {Component.get_project_directory()} && {path}")

    @staticmethod
    def get_components_directory():
        if "RUNNER_COMPONENTS_DIR" not in os.environ:
            cli_print(
                "Please set RUNNER_COMPONENTS_DIR to your components directory",
                theme_error0)
            exit(1)

        return os.environ["RUNNER_COMPONENTS_DIR"]

    @staticmethod
    def get_project_directory():
        if "RUNNER_PROJECT_DIR" not in os.environ:
            cli_print(
                "Please set RUNNER_PROJECT_DIR to your project directory",
                theme_error0)
            exit(1)

        return os.environ["RUNNER_PROJECT_DIR"]
