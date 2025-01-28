# INTERACTIVE CONTAINER

# run docker-compose from root dir
FROM ubuntu:latest

# Set non-interactive mode for apt-get to avoid blocking
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /testspace

COPY testenv/maketree.sh ./

COPY shells/tests/cli/ ./tests/

# Dependencies
RUN apt-get update && apt-get install -y \
    golang \
    git \
    curl \
    tree \
    && rm -rf /var/lib/apt/lists/*

# Install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --all

###
# fastTravelCLI install occurins in setup script
###

# Add golang binary to path
RUN echo 'export PATH=$PATH:/usr/lib/go/bin' >> ~/.bashrc

# Source setup script
RUN echo 'source ./tests/bash/interactive_setup.sh' >> ~/.bashrc



CMD ["/bin/bash"]
