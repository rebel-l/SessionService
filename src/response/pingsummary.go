package response

type PingSummary struct {
	Service,
	Storage string
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