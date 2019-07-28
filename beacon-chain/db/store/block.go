package store

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/go-ssz"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
)

func (b *BeaconDB) BeaconBlock(root []byte) (*ethpb.BeaconBlock, error) {
	return b.getBlock(root)
}

func (b *BeaconDB) HasBlock(root []byte) bool {
	if _, err := b.getBlock(root); err != nil {
		return false
	}
	return true
}

func (b *BeaconDB) SaveBeaconBlock(block *ethpb.BeaconBlock) error {
	root, err := ssz.HashTreeRoot(block)
	if err != nil {
		return err
	}
	enc, err := proto.Marshal(block)
	if err != nil {
		return err
	}
    return b.put(blocksBucket, root[:], enc)
}

func (b *BeaconDB) getBlock(root []byte) (*ethpb.BeaconBlock, error) {
	var block *ethpb.BeaconBlock
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(blocksBucket)
		enc := bkt.Get(root)
		if enc == nil {
			return fmt.Errorf("no block found for root: %#x", root)
		}
		return proto.Unmarshal(enc, block)
	}); err != nil {
		return nil, err
	}
	return block, nil
}
