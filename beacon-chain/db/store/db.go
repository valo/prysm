package store

import (
	"errors"
	"os"
	"path"
	"time"

	"github.com/boltdb/bolt"
)

// BeaconDB manages the data layer of the beacon chain implementation.
type BeaconDB struct {
	db                *bolt.DB
	DatabasePath      string
}

// NewDB returns a new Prysm DB instance.
func NewDB(dirPath string) (*BeaconDB, error) {
	if err := os.MkdirAll(dirPath, 0700); err != nil {
		return nil, err
	}
	datafile := path.Join(dirPath, "beaconchain.db")
	boltDB, err := bolt.Open(datafile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		if err == bolt.ErrTimeout {
			return nil, errors.New("cannot obtain database lock, database may be in use by another process")
		}
		return nil, err
	}

	beaconDB := &BeaconDB{db: boltDB, DatabasePath: dirPath}

	if err := beaconDB.db.Update(func(tx *bolt.Tx) error {
		return createBuckets(
			tx,
			chainMetadataBucket,
			stateBucket,
			blocksBucket,
			operationsBucket,
			validatorsBucket,
		)
	}); err != nil {
		return nil, err
	}
	return beaconDB, err
}

// Close closes the underlying database.
func (db *BeaconDB) Close() error {
	return db.db.Close()
}

// ClearDB removes the previously stored directory at the data directory.
func ClearDB(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(dirPath)
}

func createBuckets(tx *bolt.Tx, buckets ...[]byte) error {
	for _, bucket := range buckets {
		if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
			return err
		}
	}
	return nil
}
