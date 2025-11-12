#!/bin/bash -i


# Source dependency scripts
source ~/.fzf.bash
source ./maketree.sh
 
# Spoof sudo since we're root in docker
source ./tests/bash/sudo.sh

# Install fastTravelCLI
cd fastTravelCLI
source ./install/linux.sh

source ~/.bashrc

# Move back to workdir
cd /testspace/

# CLI feature tests
ft -u latest

# TODO
#  - Implement more tests
