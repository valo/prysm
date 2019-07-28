package store

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/prysmaticlabs/go-ssz"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
)

func (b *BeaconDB) Attestation(root []byte) (*ethpb.Attestation, error) {
	return b.getAttestation(root)
}

func (b *BeaconDB) HasAttestation(root []byte) bool {
	if _, err := b.getAttestation(root); err != nil {
		return false
	}
	return true
}

func (b *BeaconDB) SaveAttestation(att *ethpb.Attestation) error {
	root, err := ssz.HashTreeRoot(att)
	if err != nil {
		return err
	}
	enc, err := proto.Marshal(att)
	if err != nil {
		return err
	}
    return b.put(operationsBucket, root[:], enc)
}

func (b *BeaconDB) getAttestation(root []byte) (*ethpb.Attestation, error) {
	var att *ethpb.Attestation
	if err := b.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(operationsBucket)
		enc := bkt.Get(root)
		if enc == nil {
			return fmt.Errorf("no attestation found for root: %#x", root)
		}
		return proto.Unmarshal(enc, att)
	}); err != nil {
		return nil, err
	}
	return att, nil
}