package model

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// InitCond implements filter.InitCond
type InitCond struct {
	state *mat.VecDense
	cov   *mat.SymDense
}

// NewInitCond creates new InitCond and returns it
func NewInitCond(state mat.Vector, cov mat.Symmetric) *InitCond {
	s := &mat.VecDense{}
	s.CloneVec(state)

	c := mat.NewSymDense(cov.Symmetric(), nil)
	c.CopySym(cov)

	return &InitCond{
		state: s,
		cov:   c,
	}
}

// State returns initial state
func (c *InitCond) State() mat.Vector {
	state := mat.NewVecDense(c.state.Len(), nil)
	state.CopyVec(c.state)

	return state
}

// Cov returns initial covariance
func (c *InitCond) Cov() mat.Symmetric {
	cov := mat.NewSymDense(c.cov.Symmetric(), nil)
	cov.CopySym(c.cov)

	return cov
}

// Fall is a model of a falling ball
type Fall struct {
	// A is internal state matrix
	A *mat.Dense
	// B is control matrix
	B *mat.Dense
	// C is output state matrix
	C *mat.Dense
	// D is output control matrix
	D *mat.Dense
}

// NewFall creates a model of falling ball and returns it
func NewFall(A, B, C, D *mat.Dense) (*Fall, error) {
	return &Fall{A: A, B: B, C: C, D: D}, nil
}

// Propagate propagates internal state x of a falling ball to the next step
func (b *Fall) Propagate(x, u, q mat.Vector) (mat.Vector, error) {
	_in, _out := b.Dims()
	if u.Len() != _out {
		return nil, fmt.Errorf("Invalid input vector")
	}

	if x.Len() != _in {
		return nil, fmt.Errorf("Invalid state vector")
	}

	out := new(mat.Dense)
	out.Mul(b.A, x)

	outU := new(mat.Dense)
	outU.Mul(b.B, u)

	out.Add(out, outU)

	if q != nil && q.Len() == _in {
		out.Add(out, q)
	}

	return out.ColView(0), nil
}

// Observe observes external state of falling ball given internal state x and input u
func (b *Fall) Observe(x, u, r mat.Vector) (mat.Vector, error) {
	_in, _out := b.Dims()
	if u.Len() != _out {
		return nil, fmt.Errorf("Invalid input vector")
	}

	if x.Len() != _in {
		return nil, fmt.Errorf("Invalid state vector")
	}

	out := new(mat.Dense)
	out.Mul(b.C, x)

	outU := new(mat.Dense)
	outU.Mul(b.D, u)

	out.Add(out, outU)

	if r != nil && r.Len() == _out {
		out.Add(out, r)
	}

	return out.ColView(0), nil
}

// Dims returns input and output model dimensions
func (b *Fall) Dims() (int, int) {
	_, in := b.A.Dims()
	out, _ := b.D.Dims()

	return in, out
}