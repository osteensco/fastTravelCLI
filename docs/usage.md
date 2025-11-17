# Usage

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

# fastTravelCLI accepts args piped to it and is highly scriptable
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
