package store

import (
	"github.com/boltdb/bolt"
	"github.com/prysmaticlabs/go-ssz"
	pb "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"

)

// CanonicalState retrieves the latest, canonical state from the db.
func (b *BeaconDB) CanonicalState() (*pb.BeaconState, error) {
	var beaconState *pb.BeaconState
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc := bkt.Get(canonicalStateKey)
        return ssz.Unmarshal(enc, beaconState)
	}); err != nil {
		return nil, err
	}
	return beaconState, nil
}

// FinalizedState retrieves the finalized state from the db.
func (b *BeaconDB) FinalizedState() (*pb.BeaconState, error) {
	var beaconState *pb.BeaconState
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc := bkt.Get(finalizedStateKey)
		return ssz.Unmarshal(enc, beaconState)
	}); err != nil {
		return nil, err
	}
	return beaconState, nil
}

// JustifiedState retrieves the justified state from the db.
func (b *BeaconDB) JustifiedState() (*pb.BeaconState, error) {
	var beaconState *pb.BeaconState
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc := bkt.Get(finalizedStateKey)
		return ssz.Unmarshal(enc, beaconState)
	}); err != nil {
		return nil, err
	}
	return beaconState, nil
}

// SaveCanonicalState saves the latest, canonical state to the db.
func (b *BeaconDB) SaveCanonicalState(beaconState *pb.BeaconState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc, err := ssz.Marshal(beaconState)
		if err != nil {
			return err
		}
		return bkt.Put(canonicalStateKey, enc)
	})
}

// SaveFinalizedState saves the finalized state to the db.
func (b *BeaconDB) SaveFinalizedState(beaconState *pb.BeaconState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc, err := ssz.Marshal(beaconState)
		if err != nil {
			return err
		}
		return bkt.Put(finalizedStateKey, enc)
	})
}

// SaveJustifiedState saves the justified state to the db.
func (b *BeaconDB) SaveJustifiedState(beaconState *pb.BeaconState) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(stateBucket)
		enc, err := ssz.Marshal(beaconState)
		if err != nil {
			return err
		}
		return bkt.Put(justifiedStateKey, enc)
	})
}

