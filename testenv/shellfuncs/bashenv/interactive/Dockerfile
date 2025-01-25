# run docker-compose from root dir
FROM ubuntu:latest

# Set non-interactive mode for apt-get
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /testspace

COPY ../../../shells/tests/bashscripts/ ./

COPY ../../../shells/bash/ ./

RUN chmod +x ./exe.sh

# Update/install dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    tree \
    lua5.4 \
    && rm -rf /var/lib/apt/lists/*

# Install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --all

# Set up fzf in the shell
RUN echo 'source ~/.fzf.bash' >> ~/.bashrc

# Source scripts used to run tests
RUN echo 'source setup.sh' >> ~/.bashrc

RUN echo 'source ftmain.sh' >> ~/.bashrc

RUN echo 'source interactive_test.sh' >> ~/.bashrc



CMD ["/bin/bash"]
