package ft


// Switch for testing
var UPDATEMOCK bool = false

// user messages
const (
	InvalidDirectoryMsg       = "Provided path '%s' evaluates to '%s' which is not a valid directory. Use 'ft -ls' to see all saved destinations. \n"
	UnrecognizedKeyMsg        = "Did not recognize key or relative path '%s', use 'ft -ls' to see all saved destinations. "
	PathAlreadyExistsMsg      = "Path '%s' already exists with key '%s', overwrite key '%s'? (y/n) \n"
	AbortedOverwriteKeyMsg    = "Aborted overwriting of key '%s'. \n"
	PathOverwriteMsg          = "The value of key '%s' has been overwritten and is now '%s'. \n"
	RenamedKeyMsg             = "Renamed key '%s' to '%s' whose value is '%s'. \n"
	KeyAlreadyExistsMsg       = "Key '%s' already exists with value '%s', overwrite key '%s'? (y/n) \n"
	AddKeyMsg                 = "Added destination '%s': '%s'. \n"
	KeyDoesNotExistMsg        = "Key '%s' does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRemoveMsg           = "Are you sure you want to remove the key '%s'? (y/n) \n"
	AbortRemoveKeyMsg         = "Aborted removal of key '%s'. \n"
	RemoveKeyMsg              = "Removed '%s' destination. \n"
	RenameKeyAlreadyExistsMsg = "Key '%s' already exists, please choose something else. \n"
	RenameKeyDoesNotExistMsg  = "Cannot rename '%s', key does not exist. Run 'ft -ls' to see all keys. \n"
	VerifyRenameMsg           = "Are you sure you want to rename the key '%s' to '%s'? (y/n) \n"
	AbortRenameKeyMsg         = "Aborted renaming of key '%s' to '%s'. \n"
	IsKeyMsg                  = "Directory %s is saved to key : %s. \n"
	IsNotKeyMsg               = "No key was found for the specified path: %s. \n"
	VerifyEditMsg             = "Are you sure you want to change '%s: %s' to '%s: %s'? (y/n) \n"
	AbortEditMsg              = "Aborted replacing '%s: %s' with '%s: %s'. \n"
	PathIsNotValidDirWarn     = "Warning: Path %s is not a valid directory. \n"
)
