package arca

// DIRUD whatever
type DIRUD struct {
	Describe JSONRequestHandler
	Insert   JSONRequestHandler
	Read     JSONRequestHandler
	Update   JSONRequestHandler
	Delete   JSONRequestHandler
}
