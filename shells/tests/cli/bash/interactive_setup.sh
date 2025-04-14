#!/bin/bash

# This file must be sourced in the .bashrc to properly spoof sudo which allows fastTravelCLI install to run successfully.


# Build out the project dir and populate history
source maketree.sh

# Spoof sudo since we're root in docker
source ./tests/bash/sudo.sh

# Install fastTravelCLI
cd fastTravelCLI
source ./install/linux.sh

# This will throw an error on the line:
# source ./tests/bash/interactive_setup.sh
# This is a good thing, we don't want a recursive call to this script XD.
# We have to source .bashrc to enable the ft command. We cd back to our working directory afterwards.
source ~/.bashrc

cd /testspace/
