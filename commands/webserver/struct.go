package webserver

type GetReusltParams struct {
	PageNo   int `form:"pageNo"`
	PageSize int `form:"pageSize"`
}
