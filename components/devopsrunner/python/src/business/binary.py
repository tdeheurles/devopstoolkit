import os
from typing import Dict

from cli.Component import Component, Command
from cli.log import cli_print
from cli.parameters import key_component, key_command, indentation, key_script_name, \
    get_and_assert_additional_parameters, key
from cli.theme import theme_error0, theme_success0

component_name = "binary"
command_install = "install"
command_delete = "delete"

key_architecture = "architecture"
key_executable = "executable"
key_version = "version"
key_download_url = "download_url"

binaries_dir = f"{Component.get_project_directory()}/binaries"


def install(params: Dict[str, str], component: Component, command: Command) -> int:
    delete(params, component, command)

    all_params = get_and_assert_additional_parameters(
        initial_params=params,
        usage=component.print_component_usage,
        new_required_params=[
            {key: key_architecture},
            {key: key_executable},
            {key: key_version},
            {key: key_download_url}
        ])

    download_url = all_params[key_download_url]
    file_name = f"{all_params[key_executable]}-{all_params[key_version]}"
    s3_file_name = f"{file_name}-{all_params[key_architecture]}.tar.bz2"
    local_file_name = get_local_file_name(all_params)
    exit_code = os.system(
        f"mkdir --parents {binaries_dir}"
        f"&& cd {binaries_dir}"
        f"&& curl -O {download_url}/{s3_file_name}"
        f"&& tar --extract --bzip2 --file {s3_file_name}"
        f"&& rm {s3_file_name}"
    )

    if exit_code == 0:
        cli_print(f"installation of {local_file_name} is successful", theme_success0)
        return 0 << 8  # Don't forget to shift result by 8 bits for some python/bash reason
    else:
        cli_print(f"unable to install {local_file_name}", theme_error0)
        return 1 << 8  # Don't forget to shift result by 8 bits for some python/bash reason


def delete(params: Dict[str, str], component: Component, command: Command) -> int:
    all_params = get_and_assert_additional_parameters(
        initial_params=params,
        usage=component.print_component_usage,
        new_required_params=[
            {key: key_architecture},
            {key: key_executable},
            {key: key_version}
        ])

    local_file_name = get_local_file_name(all_params)
    if not os.path.exists(local_file_name):
        cli_print(f"{local_file_name} doesn't exist, nothing to delete", theme_success0)
        return 0 << 8

    exit_code = os.system(f"rm {local_file_name}")
    if exit_code == 0:
        cli_print(f"{local_file_name} deleted successfully", theme_success0)
    else:
        cli_print(f"unable to delete {local_file_name}", theme_error0)
    return exit_code << 8  # Don't forget to shift result by 8 bits for some python/bash reason


def get_component():
    command_base = f"{indentation}{key_script_name} {key_component} {component_name} {key_command}"
    return Component(
        name=component_name,
        definition="manipulate binaries",
        family="internal",
        commands={
            command_install: Command(
                name=command_install,
                execution=lambda params, component, command: install(params, component, command),
                help="download a binary",
                command_example=[f"{command_base} {command_install}"]
            ),
            command_delete: Command(
                name=command_delete,
                execution=lambda params, component, command: delete(params, component, command),
                help="remove local installation of a binary",
                command_example=[f"{command_base} {command_delete}"]
            )
        }
    )


def get_local_file_name(all_params: Dict[str, str]) -> str:
    return f"{binaries_dir}/{all_params[key_executable]}-{all_params[key_version]}"
