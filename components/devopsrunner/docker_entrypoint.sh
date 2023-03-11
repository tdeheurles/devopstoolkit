#!/usr/bin/env bash
set -euo pipefail

# TODO: confirm need

# We read all do.*.env files from the project root to expose them to the container
if [[ -d "${RUNNER_DOTENV_DIR}" ]];then
    if compgen -G "${RUNNER_DOTENV_DIR}/do.*.env" > /dev/null; then
        # shellcheck disable=SC2045
        for f in $(ls -A "${RUNNER_DOTENV_DIR}"/do.*.env); do
            # shellcheck disable=SC2046
            export $(grep -v '^#' "${f}" | xargs)
        done
    fi
fi

/devopsrunner $@