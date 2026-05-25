package db

import (
	"bytes"
	"fmt"

	"github.com/fluffle/golog/logging"
	"go.etcd.io/bbolt"
	"github.com/fluffle/sp0rkle/util/bson"
)

type valsBucket struct { *bbolt.Bucket }
type idxsBucket struct { *bbolt.Bucket }

func scanAll[T any](b *valsBucket) ([]T, error) {
	var results []T
	err := scanTxT(b.Bucket, func(_, v []byte) error {
		var val T
		if err := bson.Unmarshal(v, &val); err != nil {
			return err
		}
		results = append(results, val)
		return nil
	})
	return results, err
}

func scanMatch[T any](b *valsBucket, pred func(T) bool) ([]T, error) {
	var results []T
	err := scanTxT(b.Bucket, func(_, v []byte) error {
		var val T
		if err := bson.Unmarshal(v, &val); err != nil {
			return err
		}
		if pred(val) {
			results = append(results, val)
		}
		return nil
	})
	return results, err
}

func scanIndexedAll[T any](idxs *idxsBucket, vals *valsBucket) ([]T, error) {
	var results []T
	seen := map[string]bool{}
	err := scanTxT(idxs.Bucket, func(k, v []byte) error {
		if seen[string(v)] {
			return nil
		}
		seen[string(v)] = true
		data := vals.Get(v)
		if !isBson(data) {
			return nil
		}
		var val T
		if err := bson.Unmarshal(suffix(data), &val); err != nil {
			return err
		}
		results = append(results, val)
		return nil
	})
	return results, err
}

func fsckIndexPass[T Indexer](idxs *idxsBucket, vals *valsBucket) error {
	cursor := idxs.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		data := vals.Get(v)
		if !isBson(data) {
			logging.Debug("fsck: pointer %q points to non-BSON data, deleting", v)
			if delErr := cursor.Delete(); delErr != nil {
				return delErr
			}
			continue
		}
		var elem T
		if err := bson.Unmarshal(suffix(data), &elem); err != nil {
			return err
		}
		found := false
		for _, key := range elem.Indexes() {
			_, last := key.B()
			if bytes.Equal(last, k) {
				found = true
				break
			}
		}
		if !found {
			logging.Debug("fsck: obsolete index key %q, deleting", k)
			if delErr := cursor.Delete(); delErr != nil {
				return delErr
			}
		}
	}
	return nil
}

func fsckValuePass[T Indexer](idxs *idxsBucket, vals *valsBucket) error {
	cursor := vals.Cursor()
	for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
		var elem T
		if err := bson.Unmarshal(v, &elem); err != nil {
			return err
		}
		ptr := toPointer(elem)
		if !bytes.Equal(ptr, k) {
			logging.Error("fsckValuePass: key %q derived from value does not match actual key %q", ptr, k)
		}
		for _, key := range elem.Indexes() {
			elems, last := key.B()
			b := idxs.Bucket
			for _, el := range elems {
				nest := b.Bucket(el)
				if nest != nil {
					b = nest
					continue
				}
				logging.Error("fsckValuePass: index bucket %q missing, creating", el)
				var err error
				nest, err = b.CreateBucket(el)
				if err != nil {
					return fmt.Errorf("fsckValuePass: creating index bucket %q: %w", el, err)
				}
				b = nest
			}
			idxptr := b.Get(last)
			if idxptr == nil || !bytes.Equal(idxptr, ptr) {
				logging.Error("fsckValuePass: index key %q incorrect, fixing", last)
				if err := b.Put(last, ptr); err != nil {
					return fmt.Errorf("fsckValuePass: writing index pointer %q: %w", last, err)
				}
			}
		}
	}
	return nil
}

func scanTxT(b *bbolt.Bucket, scan func(k, v []byte) error) error {
	cs := []*bbolt.Cursor{b.Cursor()}
	var c *bbolt.Cursor

	for len(cs) > 0 {
		c, cs = cs[0], cs[1:]
		for k, v := c.First(); k != nil; k, v = c.Next() {
			switch {
			case v == nil:
				if nest := c.Bucket().Bucket(k); nest != nil {
					cs = append(cs, nest.Cursor())
				} else {
					logging.Error("nested bucket %q returned nil", k)
				}
			case isPointer(v):
				err := scan(k, v)
				if err == nil {
					continue
				}
				return fmt.Errorf("scan/unmarshal pointer: %w", err)
			case isBson(v):
				if err := scan(k, suffix(v)); err != nil {
					return fmt.Errorf("scan/unmarshal value: %w", err)
				}
			default:
				logging.Warn("scanTxT: unexpected data k=%q v=%q", k, v)
			}
		}
	}
	return nil
}
