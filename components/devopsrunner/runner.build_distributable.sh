#!/usr/bin/env bash
set -euo pipefail

# /!\
#   THIS COMMAND REQUIRES PYTHON THAT IS NOT INSTALLABLE WITH STD -k=dependencies
#   THIS CODE REQUIRES TO BE EXECUTED ON DEVOPSBUILDER CONTAINER
# /!\

# Build the devopsrunner as an executable in order to not require python

python_version="<% python_version %>"

(
    . "<% runner_components_dir %>/lib/bash.sh"
    cd "<% runner_project_dir %>" || exit 1

    pip install -r "<% runner_components_dir %>/<% component %>/python/src/requirements.txt"

    dist_path="<% runner_components_dir %>/<% component %>/dist"
    mkdir --parent "${dist_path}"

    work_path="$(mktemp -d)"
    pyinstaller \
        --onefile \
        --paths "<% runner_project_dir %>/binaries/python-${python_version}/site-packages" \
        --distpath "${dist_path}" \
        --workpath "${work_path}" \
        --name "<% devopsrunner_service_name %>" \
        "<% runner_components_dir %>/<% component %>/python/src/cli/do.py"
)
