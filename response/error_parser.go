package response

import "strings"

func ParseErrorResponse(errors string) map[string]string {
	errorsSplited := strings.Split(errors, ",")
	errorsParsed := map[string]string{}
	for _, error := range errorsSplited {
		errorSplited := strings.Split(error, ":")
		errorsParsed[errorSplited[0]] = errorSplited[1]
	}
	return errorsParsed
}
