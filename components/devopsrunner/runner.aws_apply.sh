#!/usr/bin/env bash
set -euo pipefail

# TODO: confirm need

# Apply the webserver to host the devopsrunner files

prompt="<% prompt %>"
disable_rollback="<% disable_rollback %>"

(
    . "<% runner_components_dir %>/lib/bash.sh"
    cd "<% runner_project_dir %>" || exit 1

    account_name="$(./do --host --component="account" --command="get_account_name")"
    if [[ "${account_name}" != "nonprod" ]];then
        theme_err "S3 DevopsRunner bucket is for now located on the non production account. Please switch login with ./login nonprod-devops"
        exit 1
    fi

    ./do --host --component="devopsrunner" \
        --command="cf_apply" \
        --disable_rollback="${disable_rollback}" \
        --stack_name="devopsrunner" \
        --aws_region="<% aws_region_main %>" \
        --prompt="${prompt}" \
        --bucket_name="<% devopsrunner_bucket_name %>"
)
