package store

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
)

// CanonicalState retrieves the latest, canonical state from the db.
func (b *BeaconDB) CanonicalState() (*pb.BeaconState, error) {
	return b.getState(stateBucket, canonicalStateKey)
}

// FinalizedState retrieves the finalized state from the db.
func (b *BeaconDB) FinalizedState() (*pb.BeaconState, error) {
	return b.getState(stateBucket, finalizedStateKey)
}

// JustifiedState retrieves the justified state from the db.
func (b *BeaconDB) JustifiedState() (*pb.BeaconState, error) {
	return b.getState(stateBucket, justifiedStateKey)
}

// SaveCanonicalState saves the latest, canonical state to the db.
func (b *BeaconDB) SaveCanonicalState(beaconState *pb.BeaconState) error {
	enc, err := proto.Marshal(beaconState)
	if err != nil {
		return err
	}
	return b.put(stateBucket, canonicalStateKey, enc)
}

// SaveFinalizedState saves the finalized state to the db.
func (b *BeaconDB) SaveFinalizedState(beaconState *pb.BeaconState) error {
	enc, err := proto.Marshal(beaconState)
	if err != nil {
		return err
	}
	return b.put(stateBucket, finalizedStateKey, enc)
}

// SaveJustifiedState saves the justified state to the db.
func (b *BeaconDB) SaveJustifiedState(beaconState *pb.BeaconState) error {
	enc, err := proto.Marshal(beaconState)
	if err != nil {
		return err
	}
	return b.put(stateBucket, justifiedStateKey, enc)
}

func (b *BeaconDB) getState(bucket []byte, key []byte) (*pb.BeaconState, error) {
	var beaconState *pb.BeaconState
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucket)
		enc := bkt.Get(key)
		if enc == nil {
			return fmt.Errorf("no item found for key: %s", key)
		}
		return proto.Unmarshal(enc, beaconState)
	}); err != nil {
		return nil, err
	}
	return beaconState, nil
}

func (b *BeaconDB) put(bucket []byte, key []byte, value []byte) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(bucket)
		return bkt.Put(key, value)
	})
}
