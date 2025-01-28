# run docker-compose from root dir
FROM ubuntu:latest

# Set non-interactive mode for apt-get to avoid blocking
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /testspace

COPY testenv/maketree.sh ./

COPY /shells/tests/cli/ ./tests/

# Dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    tree \
    && rm -rf /var/lib/apt/lists/*

# Install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --all

#####
# Install fastTravelCLI step will need to occur in each shell's test script in order for it to install properly.
#####

# Source test wrapper script
RUN chmod +x ./tests/test.sh

# Run tests
CMD ["/bin/bash", "-c", "./tests/test.sh"]
