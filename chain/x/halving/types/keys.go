package types

const (
	// ModuleName defines the module name
	ModuleName = "halving"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	HalvingParamsKey      = []byte{0x01}
	CurrentHalvingKey     = []byte{0x02}
	LastDistributionKey   = []byte{0x03}
	TotalDistributedKey   = []byte{0x04}
)