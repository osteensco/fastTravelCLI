# run docker-compose from root dir
FROM ubuntu:latest

# Set non-interactive mode for apt-get to avoid blocking
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /testspace

COPY testenv/maketree.sh ./

COPY shells/tests/cli/ ./tests/

# Dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    tree \
    && rm -rf /var/lib/apt/lists/*

# Install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --all

# Install fastTravelCLI
RUN git clone https://github.com/osteensco/fastTravelCLI.git && \
    bash ./fastTravelCLI/install/linux.sh

# Build a project dir
RUN source maketree.sh

# Run tests
CMD ["/bin/bash"]
