package store

var (
	// BoltDB buckets defined for Eth2.
	chainMetadataBucket = []byte("chain-metadata")
	stateBucket = []byte("state")
	blocksBucket = []byte("blocks")
	operationsBucket = []byte("operations")
	validatorsBucket = []byte("validators")
	// Chain metadata lookup keys.
	// TODO(#3064): Add keys...

	// State lookup keys.
	finalizedStateKey = []byte("finalized-state")
	justifiedStateKey = []byte("justified-state")
	canonicalStateKey = []byte("canonical-state")
)