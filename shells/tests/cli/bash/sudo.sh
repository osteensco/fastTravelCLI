#!/bin/bash

# Used to spoof sudo.
# Docker typically runs as root, so trying to use sudo throws an error.
# Hence, we must spoof.

sudo() {
    eval "$@"
}
