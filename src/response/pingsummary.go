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

func (ps *PingSummary) ServiceOnline() {
	ps.service = PONG
}

func (ps *PingSummary) StorageOnline() {
	ps.storage = PONG
}