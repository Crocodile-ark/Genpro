package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	// General transaction fees (40/30/30)
	KeyGeneralValidatorShare = []byte("GeneralValidatorShare")
	KeyGeneralDexShare       = []byte("GeneralDexShare")
	KeyGeneralPosShare       = []byte("GeneralPosShare")

	// LP community farming fees (30/25/25/20)
	KeyFarmingValidatorShare  = []byte("FarmingValidatorShare")
	KeyFarmingDexShare        = []byte("FarmingDexShare")
	KeyFarmingLPRewardShare   = []byte("FarmingLPRewardShare")
	KeyFarmingPosShare        = []byte("FarmingPosShare")
)

// Default parameter values for general transactions
const (
	DefaultGeneralValidatorShare = "0.40" // 40%
	DefaultGeneralDexShare       = "0.30" // 30%
	DefaultGeneralPosShare       = "0.30" // 30%
)

// Default parameter values for farming transactions
const (
	DefaultFarmingValidatorShare = "0.30" // 30%
	DefaultFarmingDexShare       = "0.25" // 25%
	DefaultFarmingLPRewardShare  = "0.25" // 25%
	DefaultFarmingPosShare       = "0.20" // 20%
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	generalValidatorShare, _ := sdk.NewDecFromStr(DefaultGeneralValidatorShare)
	generalDexShare, _ := sdk.NewDecFromStr(DefaultGeneralDexShare)
	generalPosShare, _ := sdk.NewDecFromStr(DefaultGeneralPosShare)
	
	farmingValidatorShare, _ := sdk.NewDecFromStr(DefaultFarmingValidatorShare)
	farmingDexShare, _ := sdk.NewDecFromStr(DefaultFarmingDexShare)
	farmingLPRewardShare, _ := sdk.NewDecFromStr(DefaultFarmingLPRewardShare)
	farmingPosShare, _ := sdk.NewDecFromStr(DefaultFarmingPosShare)

	return Params{
		GeneralValidatorShare: generalValidatorShare,
		GeneralDexShare:       generalDexShare,
		GeneralPosShare:       generalPosShare,
		FarmingValidatorShare: farmingValidatorShare,
		FarmingDexShare:       farmingDexShare,
		FarmingLPRewardShare:  farmingLPRewardShare,
		FarmingPosShare:       farmingPosShare,
	}
}

// ParamKeyTable the param key table for feerouter module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateShare(p.GeneralValidatorShare); err != nil {
		return fmt.Errorf("invalid general validator share: %w", err)
	}
	if err := validateShare(p.GeneralDexShare); err != nil {
		return fmt.Errorf("invalid general dex share: %w", err)
	}
	if err := validateShare(p.GeneralPosShare); err != nil {
		return fmt.Errorf("invalid general pos share: %w", err)
	}

	// Ensure general shares add up to 1.0
	generalTotal := p.GeneralValidatorShare.Add(p.GeneralDexShare).Add(p.GeneralPosShare)
	if !generalTotal.Equal(sdk.OneDec()) {
		return fmt.Errorf("general transaction shares must add up to 1.0, got %s", generalTotal.String())
	}

	if err := validateShare(p.FarmingValidatorShare); err != nil {
		return fmt.Errorf("invalid farming validator share: %w", err)
	}
	if err := validateShare(p.FarmingDexShare); err != nil {
		return fmt.Errorf("invalid farming dex share: %w", err)
	}
	if err := validateShare(p.FarmingLPRewardShare); err != nil {
		return fmt.Errorf("invalid farming LP reward share: %w", err)
	}
	if err := validateShare(p.FarmingPosShare); err != nil {
		return fmt.Errorf("invalid farming pos share: %w", err)
	}

	// Ensure farming shares add up to 1.0
	farmingTotal := p.FarmingValidatorShare.Add(p.FarmingDexShare).Add(p.FarmingLPRewardShare).Add(p.FarmingPosShare)
	if !farmingTotal.Equal(sdk.OneDec()) {
		return fmt.Errorf("farming transaction shares must add up to 1.0, got %s", farmingTotal.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements the ParamSet interface and returns the key/value pairs
// of feerouter module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyGeneralValidatorShare, &p.GeneralValidatorShare, validateShare),
		paramtypes.NewParamSetPair(KeyGeneralDexShare, &p.GeneralDexShare, validateShare),
		paramtypes.NewParamSetPair(KeyGeneralPosShare, &p.GeneralPosShare, validateShare),
		paramtypes.NewParamSetPair(KeyFarmingValidatorShare, &p.FarmingValidatorShare, validateShare),
		paramtypes.NewParamSetPair(KeyFarmingDexShare, &p.FarmingDexShare, validateShare),
		paramtypes.NewParamSetPair(KeyFarmingLPRewardShare, &p.FarmingLPRewardShare, validateShare),
		paramtypes.NewParamSetPair(KeyFarmingPosShare, &p.FarmingPosShare, validateShare),
	}
}

func validateShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("share cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("share cannot be greater than 1: %s", v)
	}

	return nil
}