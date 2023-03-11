#!/usr/bin/env bash

# TODO: confirm need

# The bootstrap is what makes the "do" versions.
# - This code let you have a few options when doing ./do --help, --verbose, ....
# - It also provides selection of --host mode or default docker mode
# - It is served so do code can download it (in S3 atm)

RUNNER_VERSION="__SEMANTIC_VERSION__"
EXIT_CODE=0

# Docker
RUNNER_DOCKER_IMAGE="${RUNNER_DOCKER_IMAGE:="ekdilymuzvlfcppzemye/devopsrunner:${RUNNER_VERSION}"}"
RUNNER_DOCKER_TTY=YES
RUNNER_DOCKER_SHARE_PROJECT="${RUNNER_DOCKER_SHARE_PROJECT:="true"}"

# Host
RUNNER_SRC=$(
    if [[ "${RUNNER_DEBUG}" == True ]]; then
        echo "${RUNNER_EXEC_DIR}/src"
    else
        echo "${RUNNER_HOME}/${RUNNER_VERSION}"
    fi
)

# Theme
INDENTATION="  "
THEME_MAIN_0="\e[36m"
THEME_ERROR_0="\e[0;91m"
NO_COLOR="\e[0m"

# COMPATIBILITY
check_do_version() {
    if [[ "${RUNNER_DO_SCRIPT_VERSION}" != 0 ]]; then
        echo -e "${THEME_ERROR_0}Your do script is outdated, please run:"
        echo -e "  curl -O https://runnerartifacts.s3.eu-west-3.amazonaws.com/master/do && chmod u+x do"
        exit 1
    fi
}

# USAGE
usage() {
    echo -e "${THEME_MAIN_0}\nUsage:${NO_COLOR}"
    echo -e "   ${0} [bootstrap option] [Runner Execution Options and Parameters]"
    echo
    echo -e "${THEME_MAIN_0}Where [Options] can be one or many from${NO_COLOR}"
    echo -e "${INDENTATION}-u | --update     update to the REQUESTED_VERSION_OR_CHANNEL defined in the ${0} script"
    echo -e "${INDENTATION}-c | --clean      clean all local files"
    echo -e "${INDENTATION}-v | --verbose    print a few information before execution"
    echo -e "${INDENTATION}-o | --host       run the execution directly on your host"
    echo -e "${INDENTATION}-t | --no-tty     remove tty for the runner docker instance"
    echo -e "${INDENTATION}-h | --help       show this help message and exit"
    echo -e "${INDENTATION}     --version    print the script version and exit"
    echo
    echo -e "${THEME_MAIN_0}Environment variables:"
    echo -e "${INDENTATION}RUNNER_CHANNEL_OR_VERSION    The version you want the runner to execute"
    echo -e "${INDENTATION}RUNNER_HOME                  The path where are stored runner files"
    echo -e "${INDENTATION}RUNNER_COMPONENTS_DIR        The path for your components (default: WORKING_DIR/components)"
    echo -e "${INDENTATION}RUNNER_EXEC_DIR              The path where to find the runner code, this is used for dev and debug"
    echo -e "${INDENTATION}RUNNER_ARTIFACT_URL          The path from where to download the artifacts"
    echo -e "${INDENTATION}RUNNER_DOTENV_DIR            The path to the directory with do.*.env files"
    echo -e "${INDENTATION}RUNNER_DOCKER_ARGS           Additional arguments for the docker run command"
    echo -e "${INDENTATION}RUNNER_DOCKER_IMAGE          The docker image used to execute code"
    echo
    echo -e "${THEME_MAIN_0}Examples:${NO_COLOR}"
    echo -e "${INDENTATION}${0}"
    echo -e "${INDENTATION}${INDENTATION}Show Runner Execution Options and Parameters"
    echo
    echo -e "${INDENTATION}${0} --help"
    echo -e "${INDENTATION}${INDENTATION}Show this help"
}

# SHOW
verbose() {
    if [[ ${VERBOSE} == YES ]]; then
        echo -e "${1}"
    fi
}
print_help() {
    if [[ ${HELP} == YES && ${ARGS} == "" ]]; then
        usage
        exit 1
    fi
}
show_variables() {
    if [[ ${VERBOSE} == YES ]]; then
        echo -e "${THEME_MAIN_0}Environment variables${NO_COLOR}"
        echo -e "${INDENTATION}RUNNER_CHANNEL_OR_VERSION     ${RUNNER_CHANNEL_OR_VERSION:-not set}"
        echo -e "${INDENTATION}RUNNER_VERSION                ${RUNNER_VERSION:-not set}"
        echo -e "${INDENTATION}RUNNER_HOME                   ${RUNNER_HOME:-not set}"
        echo -e "${INDENTATION}RUNNER_COMPONENTS_DIR         ${RUNNER_COMPONENTS_DIR:-not set}"
        echo -e "${INDENTATION}RUNNER_EXEC_DIR               ${RUNNER_EXEC_DIR:-not set}"
        echo -e "${INDENTATION}RUNNER_ARTIFACT_URL           ${RUNNER_ARTIFACT_URL:-not set}"
        echo -e "${INDENTATION}RUNNER_DEBUG                  ${RUNNER_DEBUG:-not set}"
        echo -e "${INDENTATION}RUNNER_SRC                    ${RUNNER_SRC:-not set}"
        echo -e "${INDENTATION}RUNNER_DOTENV_DIR             ${RUNNER_DOTENV_DIR:-not set}"
        echo -e "${INDENTATION}RUNNER_DOCKER_ARGS            ${RUNNER_DOCKER_ARGS:-not set}"
        echo -e "${INDENTATION}RUNNER_DO_SCRIPT_VERSION      ${RUNNER_DO_SCRIPT_VERSION:-not set}"
        echo -e "${INDENTATION}RUNNER_DOCKER_IMAGE           ${RUNNER_DOCKER_IMAGE:-not set}"
        echo -e "${INDENTATION}RUNNER_DOCKER_SHARE_PROJECT   ${RUNNER_DOCKER_IMAGE:-not set}"
        echo -e ""
        echo -e "${THEME_MAIN_0}Bootstrap internal variables${NO_COLOR}"
        echo -e "${INDENTATION}UPDATE           ${UPDATE:-NO}"
        echo -e "${INDENTATION}CLEAN            ${CLEAN:-NO}"
        echo -e "${INDENTATION}VERBOSE          ${VERBOSE:-NO}"
        echo -e "${INDENTATION}HOST_MODE        ${HOST_MODE:-NO}"
        echo -e "${INDENTATION}HELP             ${HELP:-NO}"
        echo -e "${INDENTATION}PRINT_VERSION    ${PRINT_VERSION:-NO}"
        echo -e "${INDENTATION}ARGS             ${ARGS:-NO}"
    fi
}
show_version() {
    if [[ ${PRINT_VERSION} == YES ]]; then
        echo "${RUNNER_VERSION}"
        exit 0
    fi
}

# ASSERT
assert_dependencies() {
    for c in "${@}"; do
        if [[ "$(command -v "${c}")" == "" ]]; then
            echo -e "${THEME_ERROR_0}Error, command ${c} is required to run ${0}${NO_COLOR}"
            exit 1
        fi
    done
}

# CLEAN
clean_directories() {
    echo -e "${THEME_ERROR_0}CLEAN COMMAND IS NOT IMPLEMENTED${NO_COLOR}"
    exit 1
}
clean() {
    if [[ ${CLEAN} == YES ]]; then
        clean_directories
        exit 0
    fi
}

# START
run_in_docker() {
    verbose "${THEME_MAIN_0}Run in docker${NO_COLOR}"
    assert_dependencies "docker" "printenv"

    uuid=$(uuidgen)
    container_id="stxr-devopsrunner-${uuid:0:8}"

    # UPDATE
    if [[ "${UPDATE}" == YES ]]; then
        verbose "Pulling docker image"
        docker pull "${RUNNER_DOCKER_IMAGE}"
    fi

    # BUILD COMMAND
    command="docker run"
    command+=" --name=${container_id}"
    command+=" ${RUNNER_DOCKER_ARGS:-""}"
    command+=" --interactive"

    if [[ "${RUNNER_DOCKER_SHARE_PROJECT}" == "true" ]];then
        command+=" --volume=$(pwd):/project"
    fi

    # -- ENV FILES
    if [[ -d "${RUNNER_DOTENV_DIR}" ]]; then
        if compgen -G "${RUNNER_DOTENV_DIR}/do.*.env" >/dev/null; then
            # shellcheck disable=SC2045
            for f in $(ls -A "${RUNNER_DOTENV_DIR}"/do.*.env); do
                filename=$(basename "${f}")
                command+=" --volume=${f}:/runner/${filename}"
            done
        fi
    fi

    # -- TTY
    if [[ ${RUNNER_DOCKER_TTY} == YES ]]; then
        command+=" --tty"
    fi

    command+=" --env=RUNNER_CHANNEL_OR_VERSION=${RUNNER_CHANNEL_OR_VERSION}"
    command+=" ${RUNNER_DOCKER_IMAGE}"
    command+=" ${ARGS}"

    # EXECUTE COMMAND
    verbose "${command}"
    ${command}

    EXIT_CODE=$(docker inspect "${container_id}" --format='{{.State.ExitCode}}')
    docker rm "${container_id}" >/dev/null 2>&1 || true
}
run_on_host() {
    verbose "${THEME_MAIN_0}Run on host${NO_COLOR}"
    assert_dependencies "unzip" "curl"

    if [[ -d "${RUNNER_DOTENV_DIR}" ]]; then
        if compgen -G "${RUNNER_DOTENV_DIR}/do.*.env" >/dev/null; then
            # shellcheck disable=SC2045
            for f in $(ls -A "${RUNNER_DOTENV_DIR}"/do.*.env); do
                lines=$(cat "${RUNNER_DOTENV_DIR}"/do.*.env | sed '/^[[:blank:]]*#/d;s/#.*//')
                for line in ${lines}; do
                    export $line
                done
            done
        fi
    fi

    if [[ ! -f "${RUNNER_SRC}/devopsrunner" && ! -f "${RUNNER_HOME}/src/cli/do.py" ]]; then
        verbose "Downloading devopsrunner executable because its not here and source code is not here"
        tmp_directory=${RUNNER_HOME}/tmp
        mkdir -p "${tmp_directory}"
        zip_name=${tmp_directory}/runner.${RUNNER_VERSION}.zip
        curl -s -o "${zip_name}" "${RUNNER_ARTIFACT_URL}/${RUNNER_VERSION}/runner.zip"
        mkdir -p "${RUNNER_SRC}"
        unzip -o -q "${zip_name}" -d "${RUNNER_SRC}"
        verbose "path of runner is ${RUNNER_SRC}"
        rm -r "${tmp_directory}"
    fi


    if [[ -f "${RUNNER_HOME}/src/cli/do.py" ]]; then
        verbose "Executing local source code using Python"
        python "${RUNNER_HOME}/src/cli/do.py" ${ARGS}
    else
        verbose "Executing local devopsrunner executable"
        "${RUNNER_SRC}/devopsrunner" ${ARGS}
    fi

    EXIT_CODE="${?}"
}
start_runner() {
    if [[ ${HOST_MODE} == YES ]]; then
        run_on_host
    else
        run_in_docker
    fi
}

# MAIN
main() {
    check_do_version

    # CONTROL PARAMETERS
    for i in "${@}"; do
        case ${i} in
        -u | --update)
            UPDATE=YES
            shift 1
            ;;
        -c | --clean)
            CLEAN=YES
            shift 1
            ;;
        -v | --verbose)
            VERBOSE=YES
            shift 1
            ;;
        -o | --host)
            HOST_MODE=YES
            shift 1
            ;;
        -h | --help)
            HELP=YES
            shift 1
            ;;
        -t | --no-tty)
            RUNNER_DOCKER_TTY=NO
            shift 1
            ;;
        --version)
            PRINT_VERSION=YES
            shift 1
            ;;
        *)
            ARGS="$@"
            break
            ;;
        esac
    done

    show_variables
    show_version
    print_help
    start_runner

    verbose "EXIT_CODE: ${EXIT_CODE}"
    exit ${EXIT_CODE}
}
