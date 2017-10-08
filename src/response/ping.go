package response

type Ping struct {
	Success string `json:"success"`
	Summary *PingSummary `json:"summary"`
}

func newPing(ps *PingSummary) *Ping {
	p := new(Ping)
	p.Summary = ps
	p.Success = FAILURE
	return p
}

func NewPing() *Ping {
	ps := new(PingSummary)
	return newPing(ps)
}

func (p *Ping) Notify()  {
	if p.Summary.Service == PONG && p.Summary.Storage == PONG {
		p.Success = SUCCESS
	}
}


