package domain

type constError string

func (err constError) Error() string {
	return string(err)
}

const (
	ErrGraphCycleDetected     = constError("cycle detected")
	ErrGraphModuleNotFound    = constError("Module not Found")
	ErrGraphModuleDoesntExist = constError("Module doesn't exist")
)
