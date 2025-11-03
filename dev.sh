#!/usr/bin/env bash

# Installs the dependencies of all components & starts them for development purposes.
# This script should work on a fresh clone but should be run in an environment like
# the one described in the shell.nix file.
set -euxo pipefail

echo -e "\033[1m# installing dependencies:\033[0m\n"
npm --prefix frontend/ install

echo -e "\n\n\033[1m# starting application:\033[0m\n"
overmind start -f Procfile $@
