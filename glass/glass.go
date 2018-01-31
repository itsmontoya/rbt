package glass

import (
	"github.com/Path94/atoms"
	"github.com/itsmontoya/whiskey"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidKey is returned when an invalid key is presented
	ErrInvalidKey = errors.Error("invalid key")
)

const (
	bucketPrefix = '_'
)

// New will return a new Glass
func New(dir, name string) (gp *Glass, err error) {
	var g Glass
	if g.w, err = whiskey.NewMMAP(dir, name+".wdb", 1024); err != nil {
		return
	}

	if g.s, err = whiskey.NewMMAP(dir, name+".scratch.wdb", 1024); err != nil {
		return
	}

	gp = &g
	return
}

// Glass is a database which utilizes whiskey as it's sorting algorithm
type Glass struct {
	mux atoms.RWMux

	// Master instance of whiskey
	w *whiskey.Whiskey
	// Scratch disk
	s *whiskey.Whiskey

	rtxn *Txn
	wtxn *Txn
}

// Read will return a read transaction
func (g *Glass) Read(fn TxnFn) (err error) {
	var txn Txn
	txn.r = g.w

	g.mux.Read(func() {
		err = fn(&txn)
	})

	txn.r = nil
	txn.kbuf = nil
	return
}

// Update will return an update transaction
func (g *Glass) Update(fn TxnFn) (err error) {
	var txn Txn
	txn.r = g.w
	txn.w = g.s

	g.mux.Update(func() {
		defer g.s.Reset()
		if err = fn(&txn); err != nil {
			return
		}

		err = txn.flush()
	})

	txn.r = nil
	txn.w = nil
	txn.kbuf = nil
	return
}

// Close will close an instance of glass
func (g *Glass) Close() (err error) {
	return g.w.Close()
}
