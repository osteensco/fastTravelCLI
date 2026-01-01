package data

// default value for version. This is updated at compile time.
var Version string = "development"

// GitHub endpoints for update command
var EndpointGH string = "https://api.github.com/repos/osteensco/fastTravelCLI/releases/tags/%s"
var EndpointLatestGH string = "https://api.github.com/repos/osteensco/fastTravelCLI/releases/latest"

// CLI command and resulting directory
var GitCloneCMD []string = []string{"git", "clone", "--branch", "", "https://github.com/osteensco/fastTravelCLI.git"}
var GitCloneDir string = "fastTravelCLI"
