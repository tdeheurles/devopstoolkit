import subprocess
from logging import Logger
from typing import List


def run(command: List[str], logger: Logger = None) -> [subprocess.Popen, str, str]:
    process = subprocess.Popen(
        command,
        shell=True,
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE)

    if logger is None:
        stdout, stderr = process.communicate()
        encoding = "utf-8"
        decoded_stdout = stdout.decode(encoding)
        decoded_stderr = stderr.decode(encoding)

        process.wait()
        return process, decoded_stdout, decoded_stderr

    else:
        while process.poll() is None:
            stdout, stderr = process.communicate()
            encoding = "utf-8"
            decoded_stdout = stdout.decode(encoding)
            decoded_stderr = stderr.decode(encoding)

            if decoded_stdout:
                logger.info(decoded_stdout)
            if decoded_stderr:
                logger.error(decoded_stderr)

        return process, None, None
