FROM ubuntu:22.04

# TODO: confirm need

ENV DEBIAN_FRONTEND       noninteractive
ENV RUNNER_HOME           /runner
ENV DEVOPSRUNNER_USERNAME devopsrunner
ENV RUNNER_PROJECT_DIR    /project
ENV RUNNER_COMPONENTS_DIR /project/components
ENV RUNNER_DOTENV_DIR     /project

# UTILITIES
RUN apt-get update
RUN apt-get install -y zip
RUN apt-get install -y unzip
RUN apt-get install -y less
RUN apt-get install -y curl
RUN apt-get install -y jq
RUN apt-get install -y uuid-runtime
RUN apt-get install -y wget
RUN apt-get install -y build-essential

RUN \
    aws_file="awscliv2.zip" && \
    curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" \
    --location \
    --output "${aws_file}" && \
    unzip ${aws_file} && \
    ./aws/install && \
    rm ${aws_file}

ARG SAM_VERSION=1.76.0
RUN \
    sam_file="aws-sam-cli-linux-x86_64.zip" && \
    curl "https://github.com/aws/aws-sam-cli/releases/download/v${SAM_VERSION}/${sam_file}" \
        --location \
        --output "${sam_file}" && \
    unzip ${sam_file} -d sam-installation && \
    ./sam-installation/install && \
    rm ${sam_file}

# https://go.dev/doc/install
ARG GOLANG_VERSION=1.20.1
ENV PATH="${PATH}:/usr/local/go/bin"
RUN \
    go_file="go${GOLANG_VERSION}.linux-amd64.tar.gz" && \
    curl "https://go.dev/dl/${go_file}" \
        --location \
        --output "${go_file}" && \
        tar -C /usr/local -xzf "${go_file}" && \
        rm "${go_file}"

# WGET
RUN apt-get install -y git

# RUNNER
RUN groupadd -g 1000 ${DEVOPSRUNNER_USERNAME}
RUN useradd --create-home -g ${DEVOPSRUNNER_USERNAME} -u 1000 ${DEVOPSRUNNER_USERNAME}
USER ${DEVOPSRUNNER_USERNAME}
COPY --chown=devopsrunner:devopsrunner components/devopsrunner/docker_entrypoint.sh do
COPY --chown=devopsrunner:devopsrunner components/devopsrunner/dist/devopsrunner    devopsrunner
# replace above line by a symlink probably
ARG CONTAINER_VERSION
COPY --chown=devopsrunner:devopsrunner components/devopsrunner/dist/devopsrunner    ${RUNNER_HOME}/${CONTAINER_VERSION}/devopsrunner

ENTRYPOINT ["/do"]
