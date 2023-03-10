import datetime
import os
import uuid
from typing import List, Dict

import yaml

from business import binary
from business.sam import deploy_cloudformation_stack
from cli.Component import Component, Command, key_help, key_files, key_runner_file_prefix
from cli.Usage import Usage
from cli.log import debug_print


def get_components(directories: List[str]) -> List[Component]:
    components = []
    components += get_filesystem_components(directories)
    components += get_internal_components()
    return components


def get_internal_components() -> List[Component]:
    return [
        binary.get_component()
    ]


def get_filesystem_components(directories: List[str]) -> List[Component]:
    # TODO: provide a list of directories
    components_directory = Component.get_components_directory()
    if not os.path.exists(components_directory):
        return []

    if directories == []:
        directories = [
            path
            for path in os.listdir(components_directory)
            if os.path.isdir(f"{components_directory}/{path}")
        ]

    return [
        build_filesystem_component(name=directory, directory=f"{components_directory}/{directory}")
        for directory in directories
        if "runner.ignore" not in os.listdir(f"{components_directory}/{directory}")
    ]


def build_filesystem_component(name: str, directory: str) -> Component:
    ts0 = datetime.datetime.now().timestamp()
    
    debug_print(f"\n\n=== component: {name} ===")
    usage = get_usage_from_filesystem(directory)
    debug_print(f"|_ usage:")
    debug_print(f"  |_ definition: {usage.definition}")
    debug_print(f"  |_ family:     {usage.family}")
    commands: Dict[str, Command] = {}
    for file_name in os.listdir(directory):
        tss0 = datetime.datetime.now().timestamp()
        debug_print(f"  file_name: {file_name}")
        for command_name in usage.commands.keys():
            if f"{key_runner_file_prefix}{command_name}" in file_name:
                debug_print(f"|_ command: {command_name}")

                command_usage = usage.commands[command_name]
                command_help = command_usage[key_help] if key_help in command_usage else "undefined"
                command_files = command_usage[key_files] if key_files in command_usage else []
                parsed_files = []

                # if exec file start with runner.
                file_without_prefix = Component.file_without_runner_prefix(file_name=file_name)
                execution_file = f"{str(uuid.uuid4())[0:8]}__execution__{file_without_prefix}"
                execution_template = Component.file_with_runner_prefix(file_name=execution_file)
                if file_without_prefix != file_name:
                    command_files.append(execution_template)
                    parsed_files.append(execution_template)

                debug_print(f"  |_ help:                {command_help}")
                debug_print(f"  |_ file_name:           {file_name}")
                debug_print(f"  |_ execution_template:  {execution_template}")
                debug_print(f"  |_ file_without_prefix: {file_without_prefix}")
                debug_print(f"  |_ command_files:       {command_files}")

                commands[command_name] = Command(
                    name=command_name,
                    before_parse=None if file_without_prefix == file_name
                    else lambda f=f"{directory}/{file_name}", t=f"{directory}/{execution_template}":
                        Component.execute_filesystem_before_parse(f, t),
                    execution=lambda params, component, command, path=f"{directory}/{execution_file}":
                        Component.execute_filesystem_command(path),
                    help=command_help,
                    files=command_files,
                    parsed_files=parsed_files
                )

        # Predefine commands
        if file_name == "cloudformation":
            commands["cf_apply"] = Command(
                name="cf_apply",
                help="Apply the cloudformation stack",
                execution=lambda params, component, command: deploy_cloudformation_stack(params, component, command),
                files=["cloudformation/runner.template.yaml"]
            )
        
        tss1 = datetime.datetime.now().timestamp()
        debug_print(f"build_filesystem_component: {name}.{file_name} Time consumed: {tss1 - tss0}")
    
    ts1 = datetime.datetime.now().timestamp()
    debug_print(f"build_filesystem_component: {name} Time consumed: {ts1 - ts0}")
    return Component(
        name=name,
        directory=directory,
        definition=usage.definition,
        family=usage.family,
        commands=commands)


def get_usage_from_filesystem(directory: str) -> Usage:
    usage_file = f"{directory}/runner.usage.yaml"
    definition = "No definition"
    family = "global"
    commands = {}
    if os.path.exists(usage_file):
        with open(usage_file, 'r') as stream:
            usage = yaml.safe_load(stream)
            definition = usage.get("definition", definition)
            family = usage.get("family", family)
            commands = usage.get("commands", commands)

    return Usage(definition, family, commands)
