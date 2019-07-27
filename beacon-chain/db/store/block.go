package store

import (
	"github.com/prysmaticlabs/go-ssz"

	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
)

func (b *BeaconDB) BeaconBlock(root []byte) (*ethpb.BeaconBlock, error) {
	return nil, nil
}

func (b *BeaconDB) SaveBeaconBlock(block *ethpb.BeaconBlock) error {
	root, err := ssz.HashTreeRoot(block)
	if err != nil {
		return err
	}
    return b.put(blocksBucket, root[:], block)
}
