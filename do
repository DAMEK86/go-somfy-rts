#!/usr/bin/env bash

red='\033[1;31m'
green='\033[1;32m'
normal='\033[0m'

goflags=""
if [[ "${READ_ONLY:-false}" == "true" ]]; then
    echo "Running in readonly mode"
    goflags="-mod=readonly"
fi

linter_version=1.18.0

## test
function task_test {
    go test ./... -v -count=1 $goflags && echo -e "${green}TESTS SUCCEEDED${normal}" || (echo -e "${red}!!! TESTS FAILED !!!${normal}"; exit 1)
}

function update_linter {
    pushd /tmp > /dev/null # don't install dependencies of golangci-lint in current module
    GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint@v"${linter_version}" 2>&1
    popd >/dev/null
}

## lint
function task_lint {
    update_linter
    lint_path="${GOPATH:-${HOME}/go}/bin/golangci-lint"
    "${lint_path}" run 1>&2
}

## go-fmt
function task_go_fmt {
    go fmt ./...
}

function task_build {
  os="linux"
  arch="arm"
  version=6
  GOOS=${os} GOARCH=${arch} GOARM=${version} go build -a ${goflags} -ldflags="-s -w" -o "./cli" cmd/main.go
}

function task_usage {
    echo "Usage: $0"
    sed -n 's/^##//p' <$0 | column -t -s ':' |  sed -E $'s/^/\t/'
}

CMD=${1:-}
shift || true
RESOLVED_COMMAND=$(echo "task_"$CMD | sed 's/-/_/g')
if [ "$(LC_ALL=C type -t $RESOLVED_COMMAND)" == "function" ]; then
    pushd $(dirname "${BASH_SOURCE[0]}") >/dev/null
    $RESOLVED_COMMAND "$@"
else
    task_usage
fi