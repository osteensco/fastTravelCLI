default: bash_interactive
bash_interactive:
	docker-compose -f testenv/shellfuncs/bashenv/docker-compose.yaml build
	docker-compose -f testenv/shellfuncs/bashenv/docker-compose.yaml run --rm fasttravelcli_bash_interactive

