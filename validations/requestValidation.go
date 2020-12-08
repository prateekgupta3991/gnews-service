package validations

import "errors"

func RequestQParams(qp map[string][]string, rType string) (bool, error) {
	if (rType == "source") {
		if (qp["country"] != nil && qp["language"] != nil && qp["category"] != nil){
			// return false, 
			errors.New("Invalid params. Either country and language or category")
		}
	}
	return true, nil
}