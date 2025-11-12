

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
 
fastTravelCLI provides robust bookmarking, navigation history, useful fuzzy finders (powered by fzf), and more.

fastTravelCLI is being continuously improved, check out the issues for new features, support, and integrations in the works. 


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
Fuzzy finding features require [tree](https://github.com/Old-Man-Programmer/tree) and [fzf](https://github.com/junegunn/fzf):





<h1>Usage</h1>


```bash
# Go to a directory you would like to add a bookmark for and run 
ft -set [key]


# You can also explicitly set a key to a directory, or set multiple at once
ft -set key1=some/other/dir key2=./some/relative/path


# Travel to a location by running
ft [key]


# fastTravelCLI evaluates keys to their absolute filepath, so you can do stuff like this
ft [key]/some/subdir


# ft can replace your cd command entirely and respects CDPATH
ft relative/dir
ft ..
ft -
ft mydir


# To remove a bookmark run
ft -rm [key]


# To rename a bookmark run
ft -rn [key] [new key]


# ft allows you to visit previously visited directories in your current session
ft [


# Traverse back up your dir history using 
ft ]


# fastTravelCLI has fzf integrations
# The default behavior of ft is to pull up fzf with all your bookmarks
ft

# View immediate child directories in a given project (current project by default) in fzf
ft -f
ft -f mykey
ft -f my/project/dir

# View all child directories in a given project in fzf
ft -fa
ft -fa mykey
ft -fa my/project/dir

# You can also view your session history in fzf
ft -hist


# View all your bookmarks with 
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
# you can also specify a specific version or latest
ft -u v.0.2.92



# To see a full list of available commands run
ft -help
# or
ft -h

# You can get more detailed help with a specific command
ft -set -help 
# or 
ft -set -h
```





<h1>Contributing</h1>

PRs and feature suggestions are welcome. I originally made this for myself so there's not alot of structure or official-ness to this project.
