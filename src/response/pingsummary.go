package response

type PingSummary struct {
	service,
	storage string
}

func NewPingSummary() *PingSummary {
	ps := new(PingSummary)
	return ps
}

func (ps *PingSummary) Service () string {
	return ps.service
}

func (ps *PingSummary) Storage () string {
	return ps.storage
}

func (ps *PingSummary) TurnServiceOnline() {
	ps.service = PONG
}

func (ps *PingSummary) TurnStorageOnline() {
	ps.storage = PONG
}