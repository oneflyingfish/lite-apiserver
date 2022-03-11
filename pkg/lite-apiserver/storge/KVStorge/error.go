package KVStorge

type StorgeError string

func (s StorgeError) Error() string {
	return string(s)
}

const (
	ErrNotFound  StorgeError = "Item Not Found"
	AlreadyExist StorgeError = "Item Already Exist"
	NoChange     StorgeError = "No Change to Anything"
)
