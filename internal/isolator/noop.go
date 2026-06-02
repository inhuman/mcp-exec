package isolator

import "context"

// Noop is a deterministic isolator used for unit tests of surrounding logic.
// It does not execute anything; it echoes a fixed result derived from the request.
type Noop struct {
	Result Result
	Err    error
}

func NewNoop() *Noop { return &Noop{} }

func (n *Noop) Run(_ context.Context, req Request) (Result, error) {
	if n.Err != nil {
		return Result{}, n.Err
	}
	res := n.Result
	if res.Stdout == "" {
		res.Stdout = req.Code
	}
	return res, nil
}
