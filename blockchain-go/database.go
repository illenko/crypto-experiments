package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/dgraph-io/badger/v4"
)

// Database key prefixes
const (
	BlockPrefix = "block:"
	UtxoPrefix  = "utxo:"
	MetaPrefix  = "meta:"
)

type BlockchainDB interface {
	SaveBlock(block *Block) error
	LoadBlock(index int) (*Block, error)
	SaveUTXOSet(utxoSet map[string][]*UTXO) error
	LoadUTXOSet() (map[string][]*UTXO, error)
	SaveMetadata(key string, value interface{}) error
	LoadMetadata(key string) (interface{}, error)
	GetChainHeight() (int, error)
	Close() error
}

type DatabaseManager struct {
	DB      *badger.DB
	DataDir string
}

func NewDatabaseManager(dataDir string, port int) (*DatabaseManager, error) {
	// Create path: dataDir/node-{port}/badger/
	dbPath := filepath.Join(dataDir, fmt.Sprintf("node-%d", port), "badger")

	// Ensure directory exists
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %v", err)
	}

	// Configure BadgerDB options
	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disable BadgerDB logging for cleaner output

	// Open database
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open BadgerDB: %v", err)
	}

	log.Printf("🗄️ Database opened at: %s", dbPath)

	return &DatabaseManager{
		DB:      db,
		DataDir: dbPath,
	}, nil
}

func (dm *DatabaseManager) SaveBlock(block *Block) error {
	return dm.DB.Update(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s%d", BlockPrefix, block.Index)
		data, err := json.Marshal(block)
		if err != nil {
			return fmt.Errorf("failed to marshal block: %v", err)
		}
		return txn.Set([]byte(key), data)
	})
}

func (dm *DatabaseManager) LoadBlock(index int) (*Block, error) {
	var block *Block

	err := dm.DB.View(func(txn *badger.Txn) error {
		key := fmt.Sprintf("%s%d", BlockPrefix, index)
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &block)
		})
	})

	return block, err
}

func (dm *DatabaseManager) SaveUTXOSet(utxoSet map[string][]*UTXO) error {
	return dm.DB.Update(func(txn *badger.Txn) error {
		for address, utxos := range utxoSet {
			key := fmt.Sprintf("%s%s", UtxoPrefix, address)
			data, err := json.Marshal(utxos)
			if err != nil {
				return fmt.Errorf("failed to marshal UTXOs for address %s: %v", address, err)
			}
			if err := txn.Set([]byte(key), data); err != nil {
				return err
			}
		}
		return nil
	})
}

func (dm *DatabaseManager) LoadUTXOSet() (map[string][]*UTXO, error) {
	utxoSet := make(map[string][]*UTXO)

	err := dm.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(UtxoPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			key := string(item.Key())
			address := key[len(UtxoPrefix):] // Remove prefix to get address

			err := item.Value(func(val []byte) error {
				var utxos []*UTXO
				if err := json.Unmarshal(val, &utxos); err != nil {
					return fmt.Errorf("failed to unmarshal UTXOs for address %s: %v", address, err)
				}
				utxoSet[address] = utxos
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return utxoSet, err
}

func (dm *DatabaseManager) SaveMetadata(key string, value interface{}) error {
	return dm.DB.Update(func(txn *badger.Txn) error {
		metaKey := fmt.Sprintf("%s%s", MetaPrefix, key)
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata %s: %v", key, err)
		}
		return txn.Set([]byte(metaKey), data)
	})
}

func (dm *DatabaseManager) LoadMetadata(key string) (interface{}, error) {
	var value interface{}

	err := dm.DB.View(func(txn *badger.Txn) error {
		metaKey := fmt.Sprintf("%s%s", MetaPrefix, key)
		item, err := txn.Get([]byte(metaKey))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &value)
		})
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (dm *DatabaseManager) GetChainHeight() (int, error) {
	height, err := dm.LoadMetadata("chain_height")
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return -1, nil // No chain exists yet
		}
		return -1, err
	}

	// Handle different number types from JSON unmarshaling
	switch v := height.(type) {
	case float64:
		return int(v), nil
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return -1, fmt.Errorf("invalid chain height type: %T", v)
	}
}

func (dm *DatabaseManager) Close() error {
	log.Printf("🔒 Closing database at: %s", dm.DataDir)
	return dm.DB.Close()
}
