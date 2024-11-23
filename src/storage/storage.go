package storage

import (
	"path/filepath"

	"github.com/vdbulcke/confkv/src/assert"
	"github.com/vdbulcke/confkv/src/logger"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
)

const (
	DefaultFileName = "configkv.db"
)

type OptionFunc func(*KVStore)

func WithDBName(name string) OptionFunc {
	return func(kv *KVStore) {
		if name != "" {
			kv.DBName = name

		}
	}
}
func WithLogger(l *zap.Logger) OptionFunc {
	return func(kv *KVStore) {
		kv.logger = l
	}
}

type KVStore struct {
	DBDir  string
	DBName string

	db     *bbolt.DB
	logger *zap.Logger
}

func NewKVStore(dbdir string, options ...OptionFunc) (*KVStore, error) {

	kv := &KVStore{
		DBDir:  dbdir,
		DBName: DefaultFileName,
	}

	for _, fn := range options {
		fn(kv)
	}

	if kv.logger == nil {
		kv.logger = logger.GetLogger(false)
	}

	db, err := bbolt.Open(
		filepath.Join(kv.DBDir, kv.DBName),
		0600,
		nil, // use default options
	)
	if err != nil {
		return nil, err
	}

	kv.db = db

	return kv, nil
}

func (kv *KVStore) Close() error {
	assert.NotNil(kv.db, assert.Panic, "kb.db is nil")

	return kv.db.Close()
}
