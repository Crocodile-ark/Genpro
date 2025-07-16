package app

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authvesting "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	halvingtypes "github.com/Crocodile-ark/gxrchaind/x/halving/types"
	feeroutertypes "github.com/Crocodile-ark/gxrchaind/x/feerouter/types"
)

// GXR Total Supply: 85,000,000 GXR = 8,500,000,000,000,000 ugen (8 decimals)
const (
	TotalSupplyGXR = 85_000_000 // 85 million GXR
	DecimalPlaces  = 8
	UgenPerGXR     = 100_000_000 // 1 GXR = 100,000,000 ugen
)

// GXR Supply Allocations according to specification
var (
	// Total supply in ugen
	TotalSupplyUgen = sdk.NewInt(int64(TotalSupplyGXR * UgenPerGXR))

	// Allocations in GXR (will be converted to ugen)
	AirdropFarmingGXR   = 17_000_000 // 20% - Airdrop & Farming
	DeveloperCoreGXR    = 5_950_000  // 7% - Developer Core (vesting 5 years)
	TimIntiGXR          = 5_950_000  // 7% - Tim Inti (3 orang) (soft vesting 3 years)
	LPMarketGXR         = 8_500_000  // 10% - LP & Market
	GrantGXR            = 8_500_000  // 10% - Grant (3-7 pihak)
	PoolStakingGXR      = 8_500_000  // 10% - Pool Staking (PoS)
	HalvingFundGXR      = 21_250_000 // 25% - Halving Fund
	CadanganEkspansiGXR = 8_500_000  // 10% - Cadangan/Ekspansi
	ValidatorAwalGXR    = 850_000    // 1% - Validator Awal (30 validators)
)

// GXRGenesisAllocation represents a genesis allocation
type GXRGenesisAllocation struct {
	Address     string
	Amount      sdk.Coin
	VestingType string
	VestingEnd  int64
	Description string
}

// CreateGXRGenesisAllocations creates the genesis allocations according to GXR specification
func CreateGXRGenesisAllocations(genesisTime time.Time) []GXRGenesisAllocation {
	allocations := []GXRGenesisAllocation{}

	// Convert GXR amounts to ugen
	toUgen := func(gxrAmount int64) sdk.Coin {
		return sdk.NewCoin("ugen", sdk.NewInt(gxrAmount*UgenPerGXR))
	}

	// Airdrop & Farming - distributed via Telegram bot farming (no vesting)
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1airdrop0000000000000000000000000000000000", // Placeholder address
		Amount:      toUgen(AirdropFarmingGXR),
		VestingType: "none",
		Description: "Airdrop & Farming allocation via Telegram bot",
	})

	// Developer Core - 5 year hard vesting, 10% unlock every 6 months
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1devcore0000000000000000000000000000000000", // Placeholder address
		Amount:      toUgen(DeveloperCoreGXR),
		VestingType: "continuous",
		VestingEnd:  genesisTime.Add(5 * 365 * 24 * time.Hour).Unix(), // 5 years
		Description: "Developer Core with 5-year hard vesting",
	})

	// Tim Inti (3 orang) - 3 year soft vesting
	// Split: 3% / 2% / 2%
	timIntiAmounts := []int64{
		int64(float64(TimIntiGXR) * 0.42857), // 3% / 7% = ~42.857%
		int64(float64(TimIntiGXR) * 0.28571), // 2% / 7% = ~28.571%
		int64(float64(TimIntiGXR) * 0.28571), // 2% / 7% = ~28.571%
	}
	timIntiAddresses := []string{
		"gxr1timinti1000000000000000000000000000000000", // Team member 1 (3%)
		"gxr1timinti2000000000000000000000000000000000", // Team member 2 (2%)
		"gxr1timinti3000000000000000000000000000000000", // Team member 3 (2%)
	}
	
	for i, addr := range timIntiAddresses {
		allocations = append(allocations, GXRGenesisAllocation{
			Address:     addr,
			Amount:      toUgen(timIntiAmounts[i]),
			VestingType: "continuous",
			VestingEnd:  genesisTime.Add(3 * 365 * 24 * time.Hour).Unix(), // 3 years
			Description: fmt.Sprintf("Tim Inti member %d with 3-year soft vesting", i+1),
		})
	}

	// LP & Market - initial liquidity
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1lpmarket000000000000000000000000000000000", // Placeholder address
		Amount:      toUgen(LPMarketGXR),
		VestingType: "none",
		Description: "LP & Market initial liquidity",
	})

	// Grant (3-7 pihak) - collaboration grants
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1grant00000000000000000000000000000000000", // Placeholder address
		Amount:      toUgen(GrantGXR),
		VestingType: "none",
		Description: "Grants for project and collaboration partners",
	})

	// Pool Staking (PoS) - delegator rewards
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1poolstaking00000000000000000000000000000", // Placeholder address
		Amount:      toUgen(PoolStakingGXR),
		VestingType: "none",
		Description: "PoS Pool for delegator rewards",
	})

	// Halving Fund - managed by halving module
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     authtypes.NewModuleAddress(halvingtypes.ModuleName).String(),
		Amount:      toUgen(HalvingFundGXR),
		VestingType: "none",
		Description: "Halving Fund for 5-year cycle rewards",
	})

	// Cadangan/Ekspansi - emergency and development fund
	allocations = append(allocations, GXRGenesisAllocation{
		Address:     "gxr1cadangan000000000000000000000000000000000", // Placeholder address
		Amount:      toUgen(CadanganEkspansiGXR),
		VestingType: "none",
		Description: "Emergency and ecosystem development fund",
	})

	// Validator Awal (30 validators) - early validator bonus
	// Split equally among 30 validators: 0.5% year 1, 0.5% year 2 (if active >20 days/month)
	validatorAmount := ValidatorAwalGXR / 30 // Per validator
	for i := 0; i < 30; i++ {
		allocations = append(allocations, GXRGenesisAllocation{
			Address:     fmt.Sprintf("gxr1validator%02d000000000000000000000000000", i+1),
			Amount:      toUgen(validatorAmount),
			VestingType: "continuous",
			VestingEnd:  genesisTime.Add(2 * 365 * 24 * time.Hour).Unix(), // 2 years
			Description: fmt.Sprintf("Early validator %d bonus allocation", i+1),
		})
	}

	return allocations
}

// SetupGXRGenesis configures the genesis state with GXR allocations
func SetupGXRGenesis(cdc codec.JSONCodec, genesisState GenesisState, genesisTime time.Time) GenesisState {
	// Get allocations
	allocations := CreateGXRGenesisAllocations(genesisTime)

	// Setup Auth genesis state
	var authGenState authtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[authtypes.ModuleName], &authGenState)

	// Setup Bank genesis state
	var bankGenState banktypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenState)

	// Clear existing balances and supply
	bankGenState.Balances = []banktypes.Balance{}
	bankGenState.Supply = sdk.NewCoins()

	// Add accounts and balances
	for _, alloc := range allocations {
		// Create account
		addr, err := sdk.AccAddressFromBech32(alloc.Address)
		if err != nil {
			// For placeholder addresses, skip account creation
			continue
		}

		var account authtypes.GenesisAccount
		if alloc.VestingType == "continuous" && alloc.VestingEnd > 0 {
			// Create vesting account
			baseAccount := authtypes.NewBaseAccount(addr, nil, 0, 0)
			vestingAccount := authvesting.NewContinuousVestingAccount(
				baseAccount,
				sdk.NewCoins(alloc.Amount),
				genesisTime.Unix(),
				alloc.VestingEnd,
			)
			account = vestingAccount
		} else {
			// Create regular account
			account = authtypes.NewBaseAccount(addr, nil, 0, 0)
		}

		authGenState.Accounts = append(authGenState.Accounts, account)

		// Add balance
		balance := banktypes.Balance{
			Address: alloc.Address,
			Coins:   sdk.NewCoins(alloc.Amount),
		}
		bankGenState.Balances = append(bankGenState.Balances, balance)

		// Add to total supply
		bankGenState.Supply = bankGenState.Supply.Add(alloc.Amount)
	}

	// Validate total supply
	expectedSupply := sdk.NewCoin("ugen", TotalSupplyUgen)
	if !bankGenState.Supply.IsEqual(sdk.NewCoins(expectedSupply)) {
		panic(fmt.Sprintf("Total supply mismatch: expected %s, got %s", expectedSupply, bankGenState.Supply))
	}

	// Setup Staking genesis to use ugen
	var stakingGenState stakingtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenState)
	stakingGenState.Params.BondDenom = "ugen"
	stakingGenState.Params.MaxValidators = 85

	// Setup Halving genesis
	var halvingGenState halvingtypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[halvingtypes.ModuleName], &halvingGenState)
	halvingGenState.HalvingInfo.CycleStartTime = genesisTime.Unix()

	// Setup FeeRouter genesis
	var feerouterGenState feeroutertypes.GenesisState
	cdc.MustUnmarshalJSON(genesisState[feeroutertypes.ModuleName], &feerouterGenState)

	// Marshal back to genesis state
	genesisState[authtypes.ModuleName] = cdc.MustMarshalJSON(&authGenState)
	genesisState[banktypes.ModuleName] = cdc.MustMarshalJSON(&bankGenState)
	genesisState[stakingtypes.ModuleName] = cdc.MustMarshalJSON(&stakingGenState)
	genesisState[halvingtypes.ModuleName] = cdc.MustMarshalJSON(&halvingGenState)
	genesisState[feeroutertypes.ModuleName] = cdc.MustMarshalJSON(&feerouterGenState)

	return genesisState
}