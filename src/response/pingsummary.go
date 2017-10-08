package response

type PingSummary struct {
	Service string `json:"service"`
	Storage string `json:"storage"`
}

func NewPingSummary() *PingSummary {
	ps := new(PingSummary)
	return ps
}

func (ps *PingSummary) TurnServiceOnline() {
	ps.Service = PONG
}

func (ps *PingSummary) TurnStorageOnline() {
	ps.Storage = PONG
}