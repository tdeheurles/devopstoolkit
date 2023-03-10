#!/usr/bin/env bash
set -euo pipefail

# This script can be used to locally build devopsrunner and execute command

this_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
project_root="${this_directory}/../.."
(
    cd "${project_root}" || exit 1
    # . "${project_root}/components/lib/bash.sh"
    # for f in $(ls -A "${project_root}/do.config.env"); do
    #     # shellcheck disable=SC2046
    #     export $(grep -v '^#' "${f}" | xargs)
    # done

    local_image_name_and_tag="localdevopsrunner:0"

    docker build \
        -t ${local_image_name_and_tag} \
        --build-arg="CONTAINER_VERSION=0" \
        -f ${this_directory}/Dockerfile \
        "${project_root}/."

    RUNNER_DOCKER_IMAGE="${local_image_name_and_tag}" \
        ./do $@
)
