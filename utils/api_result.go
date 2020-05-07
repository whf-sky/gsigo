package utils

func NewApiResult() *ApiResult {
	return &ApiResult{
		result: map[string]interface{}{},
	}
}

type ApiResult struct {
	code int
	result map[string]interface{}
}

func (a *ApiResult) SetCode(code int) map[string]interface{} {
	a.code = code
	return a.result
}

func (a *ApiResult) SetSuccess(data interface{}) map[string]interface{} {
	a.setResult("成功", data, 200)
	return a.result
}

func (a *ApiResult) SetError(msg string, code int) map[string]interface{} {
	a.setResult(msg, nil, code)
	return a.result
}

func (a *ApiResult) setResult(msg string, data interface{}, code int) map[string]interface{} {
	a.result["code"] = code
	a.result["data"] = data
	a.result["msg"] = msg
	return a.result
}