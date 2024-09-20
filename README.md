

<!--- <h1>fastTravelCLI</h1> --->
```
     __           _  _____                     _   ___   __   _____ - -  -  -   -   -
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ - -  -  -   -   -
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/  - -  -   -   -
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  - -  -  -   -   -
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/ - -  -  -   -   -
                                                                    
```

<h3>CD to a directory super fast</h4>

Typing out a long directory path can be a bit painful. Most tools that solve this problem give you less to type but still involve additional steps afterwards. Instead of tracking frequently visited directories, using a fuzzy search, or even a ML model to improve the cd experience, fastTravelCLI allows you to directly save a destination as a key-value pair, then you can fast travel to that location by just using the key you set. 
<br></br>
fastTravelCLI takes an experience like this: 
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
<br></br>
<br></br>

<h1>Usage</h1>
<br></br>

```rust
// Go to a directory you would like to set a waypoint for and run 

ft -set [key]


// Travel to that location by running

ft [key]


// You can also travel to a subdirectory of a key

ft [key]/some/subdir


// ft can also replace your cd command entirely

ft relative/path
ft ..
ft -
ft ./path


// You can also visit previously visited directories

ft [


// Traverse back up your dir history using 

ft ]


// View your saved locations with 

ft -ls


// To remove a location run

ft -rm [key]


// To rename a location run

ft -rn [key] [new key]


// And of course if you need help just run

ft -help
ft -h
```
<br></br>
<br></br>

<h1>Installation</h1>
<br></br>

Currently available for mac/linux OS and bash/zsh shells. May work in more shell environments but not guarenteed.

<br></br>

Requires go version >= 1.22.0

<br></br>

Clone the repo, cd into it, and run the following based on your OS - 

```
install/linux.sh
```
```
install/mac.sh
```


If you experience issues trying to run the install script as shown, you may need to prefix it with ```bash install/[myOS].sh```

<br></br>
<br></br>

<h1>Contributing</h1>
<br></br>

PRs and suggestions are welcome. I originally made this for myself and took it as a bit of a learning opportunity with some of the features I wanted to add, so there's not alot of structure or official-ness to this. This is why the install steps are just 'clone and run the shell script' instead of being available on Homebrew and other package managers. If for some reason this starts to gain a bit of interest then I'd be more than happy to implement these sort of things or accept PRs that do.
