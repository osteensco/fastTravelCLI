```
     __           _  _____                     _   ___   __   _____ - -  -  -   -   -
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ - -  -  -   -   -
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/  - -  -   -   -
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  - -  -  -   -   -
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/ - -  -  -   -   -

```

# A better CLI navigation experience

fastTravelCLI is a fast, lightweight, and feature rich CD command replacement.

fastTravelCLI provides robust bookmarking, navigation history, useful fuzzy finders (powered by fzf), and more.

fastTravelCLI is being continuously improved, check out the issues for new features, support, and integrations in the works.

## Disclaimers

Currently available for Unix-like OS and bash/zsh shells. May work in more shell environments but not guaranteed.

Compiles using go version >= 1.20.0, may work with older versions but not guaranteed.

Fuzzy finding features require [tree](https://github.com/Old-Man-Programmer/tree) and [fzf](https://github.com/junegunn/fzf):

## Contributing

PRs and feature suggestions are welcome. I originally made this for myself so there's not alot of structure or official-ness to this project.


## Getting Started

To set up the project locally for development, clone the repo and ensure you have the following installed - 

 * Docker 
 * go 1.22+
 * python 3.11+
 * lua 5.4+
 * GNU Make 4.4+

Run all tests by running the default make command or `make all`.

### Documentation Site
To run the docs site locally add a python virtual environment to the project.
```bash
python3 -m venv venv
```
Activate the virtual environment and install mkdocs and mkdocs-material.
```bash
pip install mkdocs mkdocs-material
```
You can then serve the site.
```bash
make site
```
