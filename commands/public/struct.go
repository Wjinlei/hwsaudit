package public

type Result struct {
	Id    int      `json:"id"`
	Name  string   `json:"name"`
	Path  string   `json:"path" form:"path"`
	User  string   `json:"user" form:"user"`
	Mode  string   `json:"mode" form:"mode"`
	Facl  string   `json:"facl,omitempty" form:"facl"`
	Other []string `json:"other,omitempty" form:"other"`
}

type Unit struct {
	Id          int      `json:"id"`
	Name        string   `json:"name" form:"name"`
	State       string   `json:"state" form:"state"`
	Description string   `json:"description" form:"description"`
	Path        string   `json:"path" form:"path"`
	FormState   []string `json:"form_state,omitempty" form:"form_state"`
}
