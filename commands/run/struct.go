package run

type Result struct{
    Name string `json:"name"`
    Path string `json:"path"`
    User string `json:"user"`
    Mode string `json:"mode"`
    Acl  string `json:"acl,omitempty"`
}
