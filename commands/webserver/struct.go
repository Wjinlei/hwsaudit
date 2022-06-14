package webserver

type Page struct {
	PageNo   int `form:"pageNo"`
	PageSize int `form:"pageSize"`
}
