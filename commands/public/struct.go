package public

type Result struct {
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Path  string   `json:"path" form:"path"`
	User  string   `json:"user" form:"user"`
	Mode  string   `json:"mode" form:"mode"`
	Facl  string   `json:"facl" form:"facl"`
	Other []string `json:"other" form:"other"`
}
