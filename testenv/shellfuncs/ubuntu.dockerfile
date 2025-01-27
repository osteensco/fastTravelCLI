# run docker-compose from root dir
FROM ubuntu:latest

# Set non-interactive mode for apt-get to avoid blocking
ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /testspace

COPY shells/tests/shellfuncs/ ./tests/

COPY shells/ ./shells/

COPY testenv/maketree.sh ./

RUN chmod +x ./tests/bashscripts/exe.sh

#
RUN chmod +x ./tests/bashscripts/bash_test.sh

# #
# RUN chmod +x ./tests/zshscripts/exe.sh
#
# #
# RUN chmod +x ./tests/zshscripts/zsh_test.sh

# Dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    tree \
    lua5.4 \
    && rm -rf /var/lib/apt/lists/*

# Install fzf
RUN git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf && \
    ~/.fzf/install --all

# Source test wrapper script
RUN chmod +x ./tests/test.sh

# Run tests
CMD ["/bin/bash", "-c", "./tests/test.sh"]
