

<!--- <h1>fastTravelCLI</h1> --->
```
     __           _  _____                     _   ___   __   _____ 
    / _| ____ ___| |/__   \___  ______   _____| | / __\ / /   \_   \ 
   | |_ / _  / __| __|/ /\/  _\/ _  \ \ / / _ \ |/ /   / /     / /\/ 
   |  _| (_| \__ \ |_/ /  | | | (_| |\ V /  __/ / /___/ /___/\/ /_  
   |_|  \__._|___/\__\/   |_|  \__._| \_/ \___|_\____/\____/\____/  
                                                                    
```


<h3>CD to a directory super fast</h4>

Instead of tracking frequently visited directories which will change over time, using a fuzzy search, or even a ML model to improve the cd experience, fastTravelCLI allows you to save a destination as a key-value pair, then you can fast travel to that location by just using the key you set.
<br></br>
<br></br>

<h1>Usage:</h1>
<br></br>

```golang
// Go to a directory you would like to set a waypoint for and run 

ft set [key]


// Travel to that location by running

ft to [key]


// View your saved locations with 

ft ls


// To remove a location run

ft rm [key]


// To rename a location run

ft rn [key] [new key]


// And of course if you need help just run

ft help
```

