#!/usr/bin/env bash
set -euo pipefail

# TODO: confirm need

# READ FILE
if [[ -d "${RUNNER_DOTENV_DIR}" ]];then
    if compgen -G "${RUNNER_DOTENV_DIR}/do.*.env" > /dev/null; then
        # shellcheck disable=SC2045
        for f in $(ls -A "${RUNNER_DOTENV_DIR}"/do.*.env); do
            # shellcheck disable=SC2046
            export $(grep -v '^#' "${f}" | xargs)
        done
    fi
fi

export PYTHONPATH="/runner/src"
python "${RUNNER_HOME}/src/cli/do.py" $@