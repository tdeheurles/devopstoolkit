import os
from sys import stderr

from termcolor import cprint
from cli.theme import theme_main3, theme_error0


def cli_print(log: str, color: str = theme_main3, pretty_lvl: int = 0):

    message = log
    if pretty_lvl == 1:
        message = f"\n\n---- {log}"
    if pretty_lvl == 2:
        message = f"\n\n\n========= {log} ========="

    if color == theme_error0:
        cprint(message, color, attrs=['bold'], file=stderr)
    else:
        cprint(message, color)


def debug_print(log: str):
    if "RUNNER_VERBOSE_MODE" in os.environ and os.environ["RUNNER_VERBOSE_MODE"] == "True":
        print(log)
