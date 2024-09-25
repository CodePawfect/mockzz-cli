package model

type Endpoint struct {
	API          string
	ResponseFile string
}

type Endpoints struct {
	Endpoints []Endpoint
}

func NewModel() *Endpoints {
	return &Endpoints{}
}

func (m *Endpoints) GetEndpoints() []Endpoint {
	return m.Endpoints
}

func (m *Endpoints) AddEndpoint(e Endpoint) {
	m.Endpoints = append(m.Endpoints, e)
}
