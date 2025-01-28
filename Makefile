default: all

all: shellfuncs_ubuntu cli_ubuntu

# Interactive slot for debugging
i: cli_ubuntu_bash_interactive

shellfuncs_ubuntu_bash_interactive:
	docker-compose -f testenv/docker-compose.yaml build ft_shellfuncs_ubuntu_bash_interactive
	docker-compose -f testenv/docker-compose.yaml run --rm ft_shellfuncs_ubuntu_bash_interactive 

shellfuncs_ubuntu:
	docker-compose -f testenv/docker-compose.yaml up ft_shellfuncs_ubuntu --build

cli_ubuntu_bash_interactive:
	docker-compose -f testenv/docker-compose.yaml build ft_cli_ubuntu_bash_interactive
	docker-compose -f testenv/docker-compose.yaml run --rm ft_cli_ubuntu_bash_interactive

cli_ubuntu:
	docker-compose -f test/docker-compose.yaml up ft_cli_ubuntu --build
