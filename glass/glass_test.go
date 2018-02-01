package glass

import (
	"os"
	"testing"
	"time"

	"github.com/boltDB/bolt"
	"github.com/itsmontoya/whiskey/testUtils"
	"github.com/missionMeteora/journaler"
	"github.com/missionMeteora/toolkit/errors"
)

var (
	testSortedList  = testUtils.GetSorted(100000)
	testReverseList = testUtils.GetReverse(10000)
	testRandomList  = testUtils.GetRand(10000)

	testSortedListStr  = testUtils.GetStrSlice(testSortedList)
	testReverseListStr = testUtils.GetStrSlice(testReverseList)
	testRandomListStr  = testUtils.GetStrSlice(testRandomList)

	testVal     []byte
	testBktName = []byte("testbkt")
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

		bkt.Put([]byte("name"), []byte("Josh"))
		bkt.Put([]byte("1"), []byte("1"))
		bkt.Put([]byte("2"), []byte("2"))
		bkt.Put([]byte("3"), []byte("3"))
		bkt.Put([]byte("4"), []byte("4"))
		bkt.Put([]byte("5"), []byte("5"))
		bkt.Put([]byte("6"), []byte("6"))
		bkt.Put([]byte("7"), []byte("7"))
		bkt.Put([]byte("8"), []byte("8"))
		bkt.Put([]byte("9"), []byte("9"))

		var val []byte
		if val, err = bkt.Get([]byte("name")); err != nil {
			return
		}

		journaler.Debug("Checking: %v\n", string(val))
		return
	}); err != nil {
		t.Fatal(err)
	}

	if err = g.Read(func(txn *Txn) (err error) {
		var bkt *Bucket
		if bkt = txn.Bucket([]byte("basic")); bkt == nil {
			return errors.Error("bucket doesn't exist")
		}

		if err = bkt.Put([]byte("name"), []byte("Josh")); err == nil {
			return errors.Error("expected error, received nil")
		}

		var val []byte
		if val, err = bkt.Get([]byte("name")); err != nil {
			return
		}

		journaler.Debug("Checking: %v\n", string(val))
		return
	}); err != nil {
		t.Fatal(err)
	}

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

		val, _ := bkt.Get([]byte("name"))
		journaler.Debug("Checking: %v\n", string(val))
		return
	}); err != nil {
		t.Fatal(err)
	}
}

func BenchmarkWhiskeyGet(b *testing.B) {
	var (
		g   *Glass
		err error
	)

	if err = os.MkdirAll("testing", 0755); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll("testing")

	if g, err = New("testing", "benchmarks"); err != nil {
		b.Fatal(err)
	}
	defer g.Close()

	for _, kv := range testSortedListStr {

		if err = g.Update(func(txn *Txn) (err error) {
			var bkt *Bucket
			if bkt, err = txn.CreateBucket(testBktName); err != nil {
				return
			}

			if err = bkt.Put(kv.Val, kv.Val); err != nil {
				return
			}

			return
		}); err != nil {
			b.Fatal(err)
		}

	}

	time.Sleep(time.Second)
	//	return
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range testSortedListStr {
			g.Read(func(txn *Txn) (err error) {
				bkt := txn.Bucket(testBktName)
				testVal, err = bkt.Get(kv.Val)
				return
			})
		}
	}

	b.ReportAllocs()
}

func BenchmarkWhiskeyPut(b *testing.B) {
	var (
		g   *Glass
		err error
	)

	if err = os.MkdirAll("testing", 0755); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll("testing")

	if g, err = New("testing", "benchmarks"); err != nil {
		b.Fatal(err)
	}
	defer g.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range testSortedListStr {
			g.Update(func(txn *Txn) (err error) {
				var bkt *Bucket
				if bkt, err = txn.CreateBucket(testBktName); err != nil {
					return
				}

				return bkt.Put(kv.Val, kv.Val)
			})
		}
	}

	b.ReportAllocs()
}

func BenchmarkBoltGet(b *testing.B) {
	var (
		db  *bolt.DB
		err error
	)

	if err = os.MkdirAll("testing", 0755); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll("testing")

	if db, err = bolt.Open("testing/benchmarks.bdb", 0644, nil); err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	if err = db.Update(func(txn *bolt.Tx) (err error) {
		var bkt *bolt.Bucket
		if bkt, err = txn.CreateBucket(testBktName); err != nil && err != bolt.ErrBucketExists {
			return
		}

		for _, kv := range testSortedListStr {
			if err = bkt.Put(kv.Val, kv.Val); err != nil {
				return
			}
		}

		return
	}); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range testSortedListStr {
			db.View(func(txn *bolt.Tx) (err error) {
				bkt := txn.Bucket(testBktName)
				testVal = bkt.Get(kv.Val)
				return
			})
		}
	}

	b.ReportAllocs()
}

func BenchmarkBoltPut(b *testing.B) {
	var (
		db  *bolt.DB
		err error
	)

	if err = os.MkdirAll("testing", 0755); err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll("testing")

	if db, err = bolt.Open("testing/benchmarks.bdb", 0644, nil); err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, kv := range testSortedListStr {
			db.Update(func(txn *bolt.Tx) (err error) {
				var bkt *bolt.Bucket
				if bkt, err = txn.CreateBucket(testBktName); err != nil {
					return
				}

				return bkt.Put(kv.Val, kv.Val)
			})
		}
	}

	b.ReportAllocs()
}
