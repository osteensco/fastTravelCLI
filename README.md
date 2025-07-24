

<!--- <h1>fastTravelCLI</h1> --->
```
     __           _  _____                     _   ___   __   _____ - -  -  -   -   -
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ - -  -  -   -   -
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/  - -  -   -   -
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  - -  -  -   -   -
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/ - -  -  -   -   -
                                                                    
```

<h1>A better CLI navigation experience</h1>

fastTravelCLI is a fast, lightweight, and feature rich CD command replacement.

Typing out a long directory path can be a bit painful. 
fastTravelCLI provides robust bookmarking capabilities, navigation history and more.


fastTravelCLI is being continuously improved, check out the issues for new features, support, and integrations in the works. 
Specifically, a fallback matching algorithm is in development for users that want a more automated experience.


<h1>Installation</h1>

Clone the repo, cd into it, and run the following based on your OS - 

```
bash install/linux.sh
```
```
bash install/mac.sh
```

<h2>Disclaimers</h2>
Currently available for Unix-like OS and bash/zsh shells. May work in more shell environments but not guaranteed.
<br><br>
Compiles using go version >= 1.20.0, may work with older versions but not guaranteed.
<br><br>
Some features have dependencies:

- `-hist` command requires tree and [fzf](https://github.com/junegunn/fzf).




<h1>Usage</h1>


```bash
# Go to a directory you would like to set a waypoint for and run 
ft -set [key]


# You can also explicitly set a key to a directory, or set multiple at once
ft -set key1=some/sub/dir key2=some/other/dir


# Travel to a location by running
ft [key]


# You can also travel to a subdirectory of a key
ft [key]/some/subdir


# ft can replace your cd command entirely
ft relative/dir
ft ..
ft -


# ft supports relative paths in the working directory and CDPATH
ft mydir


# To remove a location run
ft -rm [key]


# To rename a location run
ft -rn [key] [new key]


# ft allows you to visit previously visited directories
ft [


# Traverse back up your dir history using 
ft ]


# You can also view your entire session history in a fuzzy finder (powered by fzf)
ft -hist


# View your saved locations with 
ft -ls


# fastTravelCLI accepts args piped to it
echo "mykey=some/project/path" > myfile.txt
cat myfile.txt | ft -set


# if you change a directory's name on your machine, you can easily update fastTravelCLI
ft -edit my/old/dirname newdirname


# ft is easy to update to the latest release
ft -update
# or
ft -u
# you can also specify a specific version or nightly
ft -u v.0.2.92
ft -u nightly


# To see a full list of available commands run
ft -help
# or
ft -h

# You can get more detailed help with a specific command
ft -set -help 
# or 
ft -set -h
```
*NOTE:* If you run into an error opening the file `fastTravel.bin` you may need to adjust ownership of the fastTravelCLI folder.
```bash
sudo chown $USER:$USER $HOME/.local/share/fastTravelCLI
```

<h1>Common Questions</h1>

<h3>Why not just use Zoxide?</h3>

[Zoxide](https://github.com/ajeetdsouza/zoxide)  is a great tool but the implementation is complex leading to many issues that cause it to be less than ideal for many use cases.
Here are a few examples of Zoxide issues that fastTravelCLI solves gracefully or avoids entirely.

- [Ambiguous Query Matching](https://github.com/ajeetdsouza/zoxide/issues/876)
- [Lack of CDPATH support](https://github.com/ajeetdsouza/zoxide/issues/620) 
- [No session history navigation](https://github.com/ajeetdsouza/zoxide/issues/839)
- [Unable to specify query for local project](https://github.com/ajeetdsouza/zoxide/issues/863)
- [Unable to prioritize query to match child directories in cwd](https://github.com/ajeetdsouza/zoxide/issues/940)

<h1></h1>

<h3>I already use fzf for most of my in project navigation, why would I use this?</h3>

If you use [fzf](https://github.com/junegunn/fzf) to find a deeply nested directory that you access often, fastTravelCLI takes an experience like this: 
```
cd $(find * -type d | fzf)
```
```
» |
  7/7 (0)
› mystuff
  mystuff/personal
  mystuff/personal/projects/
  mystuff/personal/projects/notes
  mystuff/work
  mystuff/work/notes
  mystuff/work/docs
  mystuff/work/projects/notes
  mystuff/work/projects/
```
```
» notes|
  2/7 (0)
› mystuff/personal/projects/notes
  mystuff/work/projects/notes
```
```
» notes|
  2/7 (0)
  mystuff/personal/projects/notes
› mystuff/work/projects/notes
```
and condenses it to this:
```
ft wknotes
```


<h1>Contributing</h1>

PRs and feature suggestions are welcome. I originally made this for myself so there's not alot of structure or official-ness to this project.
