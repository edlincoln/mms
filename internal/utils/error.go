package utils

type internalServerError string

func (e internalServerError) Error() string {
	return string(e)
}

const InternalServerError = internalServerError("Internal Server Error")