default: all

all: shellfuncs

ubuntu_bash_interactive:
	docker-compose -f testenv/shellfuncs/docker-compose.yaml build fasttravelcli_ubuntu_bash_interactive
	docker-compose -f testenv/shellfuncs/docker-compose.yaml run --rm fasttravelcli_ubuntu_bash_interactive 

shellfuncs:
	docker-compose -f testenv/shellfuncs/docker-compose.yaml up fasttravelcli_shellfuncs --build
