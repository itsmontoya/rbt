package glass

import (
	"fmt"
	"os"
	"testing"
)

func TestGlass(t *testing.T) {
	var (
		g   *Glass
		err error
	)

	if err = os.MkdirAll("testing", 0755); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("testing")

	if g, err = New("testing", "data.db"); err != nil {
		t.Fatal(err)
	}
	defer g.Close()

	if err = g.Update(func(txn *Txn) (err error) {
		var bkt *Bucket
		if bkt, err = txn.CreateBucket([]byte("basic")); err != nil {
			return
		}

		bkt.w.Put([]byte("name"), []byte("Josh"))
		fmt.Printf("Checking: %v\n", string(bkt.w.Get([]byte("name"))))
		return ErrCannotWrite
	}); err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Bucket check\n")

	if err = g.Close(); err != nil {
		t.Fatal(err)
	}

	if g, err = New("testing", "data.db"); err != nil {
		t.Fatal(err)
	}

	if err = g.Update(func(txn *Txn) (err error) {
		var bkt *Bucket
		if bkt, err = txn.CreateBucket([]byte("basic")); err != nil {
			return
		}

		fmt.Printf("Checking: %v\n", string(bkt.w.Get([]byte("name"))))
		return
	}); err != nil {
		t.Fatal(err)
	}
}
