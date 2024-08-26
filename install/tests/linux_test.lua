
-- must be run from root directory (where main.go is)
-- lua install/tests/linux_test.lua

package.path = package.path .. ";install/tests/?.lua;"
require("interface")

install_path = "./install/linux.sh"

FT_DIR = "$HOME/.local/share/fastTravelCLI"
FT_EXE_PATH = FT_DIR .. "/fastTravel.exe"
FT_PATH = FT_DIR .. "/ftmain.sh"


T.set_env_vars(FT_EXE_PATH, FT_PATH, install_path)
T.main()

