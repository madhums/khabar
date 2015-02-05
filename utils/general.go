package utils

import (
	"gopkg.in/simversity/gottp.v1"
	"gopkg.in/simversity/gottp.v1/utils"
	"net/http"
)

func ConcatenateErrors(errs *[]error) string {
	var errString string
	for i := 0; i < len(*errs); i++ {
		errString += (*errs)[i].Error() + "\n"
	}
	return errString
}

func ValidateAndRaiseError(request *gottp.Request, structure interface{}) bool {
	var errs []error
	utils.ValidateStruct(&errs, structure)

	if len(errs) > 0 {
		request.Raise(gottp.HttpError{http.StatusPreconditionFailed, ConcatenateErrors(&errs)})
		return false
	}

	return true
}

func RemoveDuplicates(arr *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *arr {
		if !found[x] {
			found[x] = true
			(*arr)[j] = (*arr)[i]
			j++
		}
	}
	*arr = (*arr)[:j]
}
