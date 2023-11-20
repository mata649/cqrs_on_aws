package response

type Response interface {
	GetType() int
	GetValue() interface{}
}
type ResponseSuccessful struct {
	Value interface{}
	Type  int
}

func (r *ResponseSuccessful) GetType() int {
	return r.Type
}
func (r *ResponseSuccessful) GetValue() interface{} {
	return r.Value
}

func NewResponseSuccessful(typ int, value interface{}) *ResponseSuccessful {
	return &ResponseSuccessful{
		Value: value,
		Type:  typ,
	}
}

type ResponseFailure struct {
	Value map[string]interface{}
	Type  int
}

func (r *ResponseFailure) GetType() int {
	return r.Type
}
func (r *ResponseFailure) GetValue() interface{} {
	return r.Value
}
func NewResponseFailure(typ int, value interface{}) *ResponseFailure {
	return &ResponseFailure{
		Value: map[string]interface{}{
			"message": value,
		},
	}
}
