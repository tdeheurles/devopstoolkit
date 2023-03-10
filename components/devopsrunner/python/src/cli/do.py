import datetime
import json
import logging
import os
from typing import List, Dict

import sys
import traceback
import colorama
from sys import exit

from cli.component_loader import get_components
from cli.Component import Component
from cli.log import cli_print, debug_print
from cli.parameters import key_script_name, key_command, key_debug, key_component, keyword_component, \
    key_component_short, print_usage_options, common_parameters, get_parameters
from cli.theme import theme_main0, theme_error0, theme_main1


log_level = logging.INFO

logger = logging.getLogger()
logger.setLevel(log_level)
stream_handler = logging.StreamHandler(sys.stdout)
stream_handler.setLevel(log_level)
logger.addHandler(stream_handler)


def print_dict(dictionary):
    keys = list(dictionary.keys())
    keys.sort()
    for k in keys:
        logger.info("{0:^30}| {1}".format(k, f'{dictionary[k]}'))


def print_root_usage(components: List[Component]):
    colorama.init()
    indentation = "  "
    cli_print('\nUsage:', theme_main0)
    cli_print(
        f'{indentation}{key_script_name}'
        f' --{key_component} {keyword_component}'
        f' [Options]'
    )

    cli_print(f'\nWhere (--{key_component}|-{key_component_short}) {keyword_component} is one of:', theme_main0)

    max_length = 15
    for component in components:
        if len(component.name) > max_length:
            max_length = len(component.name)
    component_format = "{0}{1:<" + f"{max_length + 5}" + "}{2}"  # "{0:<" + f"{max_length + 5}" + "}{1:<15}"

    families = list(set([component.family for component in components]))
    families.sort()
    sorted_components = components
    sorted_components.sort(key=lambda c: c.name)
    for family in families:
        cli_print(f"{indentation}{family}", theme_main1)
        for component in sorted_components:
            if component.family == family:
                cli_print(component_format.format(indentation*2, component.name, component.definition))
        print()

    print_usage_options()


def execute_component_command(component: Component, params: Dict[str, str]) -> int:
    if key_command not in params:
        component.print_component_usage()
        cli_print(f"\n--{key_command} parameter is missing", theme_error0)
        exit(1)

    if params[key_command] not in component.commands.keys():
        component.print_component_usage()
        print()
        cli_print(
            f"Command <{params[key_command]}> doesn't exist for component {component.name}",
            theme_error0)
        exit(1)

    return component.execute_command(params)


def main():
    try:
        params = get_parameters(common_parameters)

        if key_debug in params and params[key_debug] == "true":
            os.environ["RUNNER_VERBOSE_MODE"] = "True"
            cli_print("environment", theme_main0, pretty_lvl=2)
            print_dict(os.environ)
            cli_print("parameters", theme_main0, pretty_lvl=2)
            print_dict(params)

        ts0 = datetime.datetime.now().timestamp()
        if params.get("component", "") != "":
            components = get_components([params.get("component")])
        else:
            components = get_components([])
        
        ts1 = datetime.datetime.now().timestamp()
        debug_print(f"get_filesystem_components: Time consumed: {ts1 - ts0}")

        if len(components) == 0:
            print_root_usage(components)
            cli_print(
                f"No component have been defined, please create one and retry",
                theme_error0)
            exit(1)

        if key_component not in params:
            print_root_usage(components)
            cli_print(f"\n--{key_component} parameter is missing", theme_error0)
            exit(1)

        for component in components:
            debug_print("")
            debug_print(str(component))

        for component in components:
            if component.name == params[key_component]:
                exit_code = execute_component_command(component, params)
                exit(exit_code)

        print_root_usage(components)
        cli_print(f"Component {params[key_component]} doesn't exist", theme_error0)
        exit(1)

    except Exception:
        exc_type, exc_value, exc_traceback = sys.exc_info()
        message = traceback.format_exception(exc_type, exc_value, exc_traceback)
        logger.error(json.dumps({"unhandledException": message}, indent=2))
        exit(1)


main()
