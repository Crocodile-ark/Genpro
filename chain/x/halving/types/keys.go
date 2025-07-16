package types

var (
	// Keys for store
	CurrentHalvingKey     = []byte("current_halving")
	LastDistributionKey   = []byte("last_distribution")
	ValidatorUptimeKey    = []byte("validator_uptime")
)

const (
	// ModuleName is the name of the halving module
	ModuleName = "halving"
	
	// StoreKey is the store key string for the halving module
	StoreKey = ModuleName
	
	// RouterKey is the message route for the halving module
	RouterKey = ModuleName
	
	// QuerierRoute is the querier route for the halving module
	QuerierRoute = ModuleName
)