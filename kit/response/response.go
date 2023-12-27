package response

import (
	"fmt"
	"strings"
)

type Response struct {
	value interface{}
	typ   int
}

func NewResponse(typ int, value interface{}) Response {
	return Response{
		value: value,
		typ:   typ,
	}
}

func (r Response) GetType() int {
	return r.typ
}
func (r Response) GetValue() interface{} {
	return r.value
}

func ParseErrorResponse(errors string) string {
	errorsSplited := strings.Split(errors, ",")
	errorsParsed := `{`
	for _, error := range errorsSplited {
		errorSplited := strings.Split(error, ":")
		errorsParsed = fmt.Sprintf(`"%s":"%s",`, errorSplited[0], errorSplited[1])
	}
	errorsParsed = errorsParsed[:len(errorsParsed)-1] + "}"
	return errorsParsed
}
