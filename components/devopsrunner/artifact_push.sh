#!/usr/bin/env bash
set -euo pipefail

# TODO: confirm need

# This non devopsrunner command purpose is to push the devopsrunner codebase as artifacts in s3 bucket
# It also push the bootstrap file that let us switch the version in the do for the devopsrunner
# Without that, do wouldn't be able to version/download the files

this_directory="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"
project_root="${this_directory}/../.."
(
    cd "${project_root}" || exit 1
    . "${this_directory}/../lib/bash.sh"

    build_pipeline=${BUILD_PIPELINE:?$(theme_err "Please provide environment variable BUILD_PIPELINE" && exit 1)} # CI_PIPELINE_ID
    git_branch=${GIT_BRANCH:?$(theme_err "Please provide environment variable GIT_BRANCH" && exit 1)}             # CI_COMMIT_BRANCH
    git_sha=${GIT_SHORT_SHA:?$(theme_err "Please provide environment variable GIT_SHORT_SHA" && exit 1)}    # CI_COMMIT_SHORT_SHA
    semantic_versioning=${SEMANTIC_VERSIONING:?$(theme_err "Please provide environment variable SEMANTIC_VERSIONING" && exit 1)}

    git_short_sha="${git_sha:0:7}"
    full_semantic_version="${semantic_versioning}.${git_short_sha}"

    parsed_bootstrap_file="$(mktemp)"
    sed "s/__SEMANTIC_VERSION__/${full_semantic_version}/g" "${this_directory}/bootstrap.sh" \
        > "${parsed_bootstrap_file}"

    distributable_zip="$(mktemp).zip"
    s3_base_path="s3://structure-nonprod-devopsrunner"
    (
        # `./do -k=devopsrunner -c=build_distributable` is required before
        cd "${this_directory}/dist"
        zip -r "${distributable_zip}" .
    )

    # Zip is only stored to the full semantic version that is defined in bootstrap.sh
    theme_1 "Pushing full_semantic_version"
    s3_original_bootstrap="${s3_base_path}/${full_semantic_version}/bootstrap.sh"
    s3_original_do="${s3_base_path}/${full_semantic_version}/do"
    aws s3 cp "${distributable_zip}"     "${s3_base_path}/${full_semantic_version}/runner.zip"
    aws s3 cp "${parsed_bootstrap_file}" "${s3_original_bootstrap}"
    aws s3 cp "${this_directory}/do"     "${s3_original_do}"

    # We set a version of do and bootstrap to be able to download them
    theme_1 "Duplicate do and bootstrap as pointers"
    for sub_path in "${semantic_versioning}" "${build_pipeline}" "${git_branch}" "${git_sha}"
    do
        aws s3 cp "${s3_original_bootstrap}" "${s3_base_path}/${sub_path}/bootstrap.sh"
        aws s3 cp "${s3_original_do}" "${s3_base_path}/${sub_path}/do"
    done
)
