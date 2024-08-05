
-- must be run from root directory (where main.go is)

package.path = package.path .. ";install/tests/?.lua;"
require("interface")

install_path = "./install/win.sh"

FT_EXE_PATH = "/mnt/c/fastTravelCLI/fastTravel.exe"
FT_PATH = "~/ftmain.sh"


T.set_env_vars(FT_EXE_PATH, FT_PATH, install_path)
T.main()

