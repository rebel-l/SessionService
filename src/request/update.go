package request

type Update struct {
	Id string `json:"id"`
	Data map[string]string `json:"data"`
}
