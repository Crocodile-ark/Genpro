package types

const (
	// ModuleName defines the module name
	ModuleName = "feerouter"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for fee router
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName
)

// KVStore keys
var (
	FeeRouterParamsKey = []byte{0x01}
	FeeStatsKey        = []byte{0x02}
	LPPoolsKey         = []byte{0x03}
)