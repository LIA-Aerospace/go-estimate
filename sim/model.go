package sim

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
	s.CloneFromVec(state)

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
	state.CloneFromVec(c.state)

	return state
}

// Cov returns initial covariance
func (c *InitCond) Cov() mat.Symmetric {
	cov := mat.NewSymDense(c.cov.Symmetric(), nil)
	cov.CopySym(c.cov)

	return cov
}

// BaseModel is a basic model of a dynamical system
type BaseModel struct {
	// A is internal state matrix
	A *mat.Dense
	// B is control matrix
	B *mat.Dense
	// C is output state matrix
	C *mat.Dense
	// D is output control matrix
	D *mat.Dense
}

// NewBaseModel creates a model of falling ball and returns it
func NewBaseModel(A, B, C, D *mat.Dense) (*BaseModel, error) {
	return &BaseModel{A: A, B: B, C: C, D: D}, nil
}

// Propagate propagates internal state x of a falling ball to the next step
func (b *BaseModel) Propagate(x, u, q mat.Vector) (mat.Vector, error) {
	_in, _out := b.Dims()
	if u != nil && u.Len() != _out {
		return nil, fmt.Errorf("Invalid input vector")
	}

	if x.Len() != _in {
		return nil, fmt.Errorf("Invalid state vector")
	}

	out := new(mat.Dense)
	out.Mul(b.A, x)

	if u != nil && b.B != nil {
		outU := new(mat.Dense)
		outU.Mul(b.B, u)

		out.Add(out, outU)
	}

	if q != nil && q.Len() == _in {
		out.Add(out, q)
	}

	return out.ColView(0), nil
}

// Observe observes external state of falling ball given internal state x and input u
func (b *BaseModel) Observe(x, u, r mat.Vector) (mat.Vector, error) {
	_in, _out := b.Dims()
	if u != nil && u.Len() != _out {
		return nil, fmt.Errorf("Invalid input vector")
	}

	if x.Len() != _in {
		return nil, fmt.Errorf("Invalid state vector")
	}

	out := new(mat.Dense)
	out.Mul(b.C, x)

	if u != nil && b.D != nil {
		outU := new(mat.Dense)
		outU.Mul(b.D, u)

		out.Add(out, outU)
	}

	if r != nil && r.Len() == _out {
		out.Add(out, r)
	}

	return out.ColView(0), nil
}

// Dims returns input and output model dimensions
func (b *BaseModel) Dims() (int, int) {
	_, in := b.A.Dims()
	out, _ := b.C.Dims()

	return in, out
}

// StateMatrix returns state propagation matrix
func (b *BaseModel) StateMatrix() mat.Matrix {
	m := &mat.Dense{}
	m.CloneFrom(b.A)

	return m
}

// StateCtlMatrix returns state propagation control matrix
func (b *BaseModel) StateCtlMatrix() mat.Matrix {
	m := &mat.Dense{}
	if b.B != nil {
		m.CloneFrom(b.B)
	}

	return m
}

// OutputMatrix returns observation matrix
func (b *BaseModel) OutputMatrix() mat.Matrix {
	m := &mat.Dense{}
	m.CloneFrom(b.C)

	return m
}

// OutputCtlMatrix returns observation control matrix
func (b *BaseModel) OutputCtlMatrix() mat.Matrix {
	m := &mat.Dense{}
	if b.D != nil {
		m.CloneFrom(b.D)
	}

	return m
}
