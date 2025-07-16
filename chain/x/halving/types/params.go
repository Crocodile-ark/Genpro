package types

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	KeyHalvingCycleDuration = []byte("HalvingCycleDuration")
	KeyValidatorShare       = []byte("ValidatorShare")
	KeyDelegatorShare       = []byte("DelegatorShare")
	KeyDexShare            = []byte("DexShare")
)

// Default parameter values
const (
	DefaultHalvingCycleDuration = 5 * 365 * 24 * time.Hour // 5 years
	DefaultValidatorShare       = "0.70"                   // 70%
	DefaultDelegatorShare       = "0.20"                   // 20%
	DefaultDexShare            = "0.10"                   // 10%
)

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	validatorShare, _ := sdk.NewDecFromStr(DefaultValidatorShare)
	delegatorShare, _ := sdk.NewDecFromStr(DefaultDelegatorShare)
	dexShare, _ := sdk.NewDecFromStr(DefaultDexShare)

	return Params{
		HalvingCycleDuration: DefaultHalvingCycleDuration,
		ValidatorShare:       validatorShare,
		DelegatorShare:       delegatorShare,
		DexShare:            dexShare,
	}
}

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// Validate validates the set of params
func (p Params) Validate() error {
	if err := validateHalvingCycleDuration(p.HalvingCycleDuration); err != nil {
		return err
	}
	if err := validateValidatorShare(p.ValidatorShare); err != nil {
		return err
	}
	if err := validateDelegatorShare(p.DelegatorShare); err != nil {
		return err
	}
	if err := validateDexShare(p.DexShare); err != nil {
		return err
	}

	// Ensure shares add up to 1.0
	total := p.ValidatorShare.Add(p.DelegatorShare).Add(p.DexShare)
	if !total.Equal(sdk.OneDec()) {
		return fmt.Errorf("validator, delegator, and dex shares must add up to 1.0, got %s", total.String())
	}

	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs implements the ParamSet interface and returns the key/value pairs
// of halving module's parameters.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyHalvingCycleDuration, &p.HalvingCycleDuration, validateHalvingCycleDuration),
		paramtypes.NewParamSetPair(KeyValidatorShare, &p.ValidatorShare, validateValidatorShare),
		paramtypes.NewParamSetPair(KeyDelegatorShare, &p.DelegatorShare, validateDelegatorShare),
		paramtypes.NewParamSetPair(KeyDexShare, &p.DexShare, validateDexShare),
	}
}

func validateHalvingCycleDuration(i interface{}) error {
	v, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("halving cycle duration must be positive: %d", v)
	}

	return nil
}

func validateValidatorShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("validator share cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("validator share cannot be greater than 1: %s", v)
	}

	return nil
}

func validateDelegatorShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("delegator share cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("delegator share cannot be greater than 1: %s", v)
	}

	return nil
}

func validateDexShare(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("dex share cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("dex share cannot be greater than 1: %s", v)
	}

	return nil
}