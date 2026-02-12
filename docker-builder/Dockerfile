FROM ubuntu:22.04 as base
WORKDIR /workdir

ARG sdk_nrf_branch=v2.6.1
ARG toolchain_version=v2.6.1
ARG sdk_nrf_commit

ENV DEBIAN_FRONTEND=noninteractive

SHELL [ "/bin/bash", "-euxo", "pipefail", "-c" ]

# Install dependencies
RUN <<EOT
    apt-get -y update
    apt-get -y upgrade
    apt-get -y install wget unzip clang-format make libffi7 git
    apt-get -y clean
    rm -rf /var/lib/apt/lists/*
EOT

# Install Go 1.22
RUN <<EOT
    wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.22.5.linux-amd64.tar.gz
    ln -s /usr/local/go/bin/go /usr/local/bin/go
    ln -s /usr/local/go/bin/gofmt /usr/local/bin/gofmt
    rm go1.22.5.linux-amd64.tar.gz
EOT

# Install toolchain
# Make nrfutil install in a shared location, because when used with GitHub
# Actions, the image will be launched with the home dir mounted from the local
# checkout.
ENV NRFUTIL_HOME=/usr/local/share/nrfutil
RUN <<EOT
    wget -q https://developer.nordicsemi.com/.pc-tools/nrfutil/x64-linux/nrfutil
    mv nrfutil /usr/local/bin
    chmod +x /usr/local/bin/nrfutil
    nrfutil install toolchain-manager
    nrfutil install toolchain-manager search
    nrfutil toolchain-manager install --ncs-version ${toolchain_version}
    nrfutil toolchain-manager list
    rm -f /root/ncs/downloads/*
EOT

#
# ClangFormat
#
RUN <<EOT
    wget -qO- https://raw.githubusercontent.com/nrfconnect/sdk-nrf/${sdk_nrf_branch}/.clang-format > /workdir/.clang-format
EOT

# Prepare image with a ready to use build environment
SHELL ["nrfutil","toolchain-manager","launch","/bin/bash","--","-c"]
RUN <<EOT
    west init -m https://github.com/nrfconnect/sdk-nrf --mr ${sdk_nrf_branch} .
    if [[ $sdk_nrf_commit =~ "^[a-fA-F0-9]{32}$" ]]; then
        git checkout ${sdk_nrf_commit};
    fi
    west update --narrow -o=--depth=1
EOT

ARG zigbee_home_version=develop

# Zigbee_Home
RUN <<EOT
  git clone --branch ${zigbee_home_version} https://github.com/ffenix113/zigbee_home.git /workdir/zigbee_home
  cd /workdir/zigbee_home
  go build -o zigbee_home cmd/zigbee/main.go
  mv zigbee_home /usr/local/bin
EOT

ENV ZEPHYR_BASE=/workdir/zephyr

# Launch into build environment with the passed arguments
# Currently this is not supported in GitHub Actions
# See https://github.com/actions/runner/issues/1964
ENTRYPOINT [ "nrfutil", "toolchain-manager", "launch", "--shell"]
