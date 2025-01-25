default: all

all: bash_script

bash_script_interactive:
	docker-compose -f testenv/shellfuncs/docker-compose.yaml build fasttravelcli_bash_script_interactive
	docker-compose -f testenv/shellfuncs/docker-compose.yaml run --rm fasttravelcli_bash_script_interactive

bash_script:
	docker-compose -f testenv/shellfuncs/docker-compose.yaml up fasttravelcli_bash_script --build
