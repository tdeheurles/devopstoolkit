# MAIN SCRIPT
import argparse
import os
from sys import exit
from typing import Dict, List
from termcolor import cprint
from cli.log import cli_print
from cli.theme import theme_error0, theme_main0

key_script_name = "./do"

# COMMANDS
command_deploy = "deploy"
command_apply = "apply"
command_delete = "delete"
command_get_policy = "get_policy"
available_commands = [
    command_deploy,
    command_apply,
    command_delete,
    command_get_policy,
]


# PARAMETERS
key = "key"
key_alias = "aliases"
key_debug = "debug"
key_command = "command"
key_command_short = "c"
key_s3_definition = "s3_definition"
key_catalog_payload = "catalog_payload"
key_environment = "environment_name"
key_environment_short = "e"
key_start_date = "start_date"
key_end_date = "end_date"
key_component = "component"
key_component_short = "k"
key_version = "version"
key_version_short = "v"
key_username = "username"
key_bucket_name = 'bucket_name'
key_bucket_name_short = 'b'
key_aws_region = 'aws_region'
key_stack_name = 'stack_name'
key_prompt = 'prompt'
key_disable_rollback = 'disable_rollback'
key_aws_region_short = 'r'
common_parameters = [
    {key: key_command, key_alias: key_command_short},
    {key: key_s3_definition},
    {key: "deploy"},
    {key: key_debug},
    {key: key_environment, key_alias: key_environment_short},
    {key: key_start_date},
    {key: key_end_date},
    {key: key_component, key_alias: key_component_short},
    {key: key_version, key_alias: key_version_short},
    {key: key_catalog_payload},
    {key: key_username},
    {key: key_stack_name},
]


# MISCELLANEOUS
indentation = "  "
keyword_command = "<COMMAND>"
keyword_component = "<COMPONENT>"


def assert_component(usage, params: Dict[str, str]):
    if key_component not in params and key_component_short not in params:
        usage()
        cprint(f"\n[--{key_component}|-{key_component_short}] parameter is missing", theme_error0)
        exit(1)


def standard_key_assert(usage, key_name, params: Dict[str, str], message: str = "", shortkey_name: str = None):
    if key_name not in params:
        usage()
        shortkey = f"|-{shortkey_name}" if shortkey_name is not None else ""
        cprint(f"\n--{key_name}{shortkey} parameter is missing. {message}", theme_error0)
        exit(1)


def print_usage_options():
    command_format = indentation + "{0:<25}{1:<15}"
    cli_print(f'Where [Options] can be one or many from:', theme_main0)
    cli_print(command_format.format(
        f"--{key_debug} (true|false)",
        "Print all environment variables and received parameters (env + cli) - Default (false)"))


def get_parameters_from_environment(parameters: List[Dict[str, str]]):
    parameters_from_environment = {}

    for parameter in parameters:
        v = os.environ.get(parameter[key].upper(), None)
        if v is not None:
            parameters_from_environment[parameter[key]] = v

    return parameters_from_environment


def get_parameter_from_args(parameters: List[Dict[str, str]]) -> dict:
    parser = argparse.ArgumentParser()

    for parameter in parameters:
        if key_alias in parameter:
            parser.add_argument(f'-{parameter[key_alias]}', f'--{parameter[key]}', required=False)
        else:
            parser.add_argument(f'--{parameter[key]}', required=False)

    namespace, unknown = parser.parse_known_args()

    parameters_from_args = {}
    for parameter in parameters:
        value = getattr(namespace, parameter[key])
        if value is not None:
            parameters_from_args[parameter[key]] = value

    return parameters_from_args


def get_parameters(parameters: List[Dict[str, str]]):
    params = get_parameters_from_environment(parameters)
    params.update(get_parameter_from_args(parameters))
    return params


def get_and_assert_additional_parameters(
        initial_params: Dict[str, str],
        usage, new_required_params: List[Dict[str, str]]):

    new_params = {}
    for element in new_required_params:
        for _, v in element.items():
            if v not in initial_params:
                new_params[v] = {key: v}

    all_params = {
        **initial_params,
        **get_parameters(list(new_params.values()))
    }

    assert_parameters(
        params=all_params,
        usage=usage,
        new_required_params=list(new_params.values())
    )

    return all_params


def get_additional_parameters(initial_params: Dict[str, str], new_required_params: List[Dict[str, str]]
                              ) -> Dict[str, str]:
    return {
        **initial_params,
        **get_parameters(new_required_params)
    }


def assert_parameters(params: Dict[str, str], usage, new_required_params: List[Dict[str, str]]):

    missing_parameters = []
    for param in new_required_params:
        if param[key] not in params:
            missing_parameters.append(param[key])

    if len(missing_parameters) != 0:
        usage()
        print()
        for missing_param in missing_parameters:
            cli_print(f"--{missing_param} parameter is missing", theme_error0)
        exit(1)
