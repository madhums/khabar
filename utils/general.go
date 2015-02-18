package utils

import (
	"github.com/changer/sc-notifications/db"
	"gopkg.in/simversity/gottp.v1"
	"gopkg.in/simversity/gottp.v1/utils"
	"net/http"
)

func ConcatenateErrors(errs *[]error) string {
	var errString string
	for i := 0; i < len(*errs); i++ {
		errString += (*errs)[i].Error()
		if (len(*errs) - i) > 1 {
			errString += ","
		}
	}
	return errString
}

func GetPaginationToQuery(paginator *gottp.Paginator) *db.M {
	query := make(db.M)
	if len(paginator.Ids) > 0 {
		query["_id"] = db.M{
			"$in": paginator.Ids,
		}
	}

	if len(paginator.Wkey) > 0 {
		if len(paginator.Wgt) > 0 {
			query[paginator.Wkey] = db.M{
				"$gt": paginator.Wgt,
			}
		}

		if len(paginator.Wlt) > 0 {
			query[paginator.Wkey] = db.M{
				"$lt": paginator.Wlt,
			}
		}
	}

	if paginator.Limit > 0 {
		query["limit"] = paginator.Limit
	}

	if paginator.Skip > 0 {
		query["skip"] = paginator.Skip
	}
	return &query
}

func ValidateAndRaiseError(request *gottp.Request, structure interface{}) bool {
	var errs []error
	utils.ValidateStruct(&errs, structure)

	if len(errs) > 0 {
		request.Raise(gottp.HttpError{http.StatusBadRequest, ConcatenateErrors(&errs)})
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
