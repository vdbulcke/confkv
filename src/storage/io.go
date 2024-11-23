package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/vdbulcke/confkv/src/assert"
	"go.etcd.io/bbolt"
)

func (kv *KVStore) Get(ctx context.Context, bucketName, key string) ([]byte, error) {
	assert.NotNil(kv.db, assert.Panic, "kb.db is nil")

	var value []byte
	err := kv.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return nil
		}

		v := b.Get([]byte(key))

		// copy(value, v)
		value = bytes.Clone(v)

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Get err: %w", err)
	}

	return value, nil
}

func (kv *KVStore) Put(ctx context.Context, bucket, key string, value []byte) error {
	assert.NotNil(kv.db, assert.Panic, "kb.db is nil")

	err := kv.db.Update(func(tx *bbolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			return e
		}

		return b.Put([]byte(key), value)

	})

	return err
}

func (kv *KVStore) Delete(ctx context.Context, bucket, key string) error {
	assert.NotNil(kv.db, assert.Panic, "kb.db is nil")

	err := kv.db.Update(func(tx *bbolt.Tx) error {
		b, e := tx.CreateBucketIfNotExists([]byte(bucket))
		if e != nil {
			return e
		}

		return b.Delete([]byte(key))
	})

	return err
}

func (kv *KVStore) DeleteBucket(ctx context.Context, bucket string) error {
	assert.NotNil(kv.db, assert.Panic, "kb.db is nil")

	err := kv.db.Update(func(tx *bbolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))

	})
	if err != nil && errors.Is(err, bbolt.ErrBucketNotFound) {
		return nil
	}

	return err
}
