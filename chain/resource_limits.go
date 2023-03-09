package chain

import (
	"fmt"
	"math"

	"github.com/MetalBlockchain/antelopevm/config"
	"github.com/MetalBlockchain/antelopevm/core/name"
	"github.com/MetalBlockchain/antelopevm/core/resource"
	"github.com/MetalBlockchain/antelopevm/state"
)

type ResourceLimitsManager struct {
	session *state.Session
}

func NewResourceLimitsManager(session *state.Session) *ResourceLimitsManager {
	return &ResourceLimitsManager{
		session: session,
	}
}

func (rl *ResourceLimitsManager) GetBlockCpuLimit() uint32 {
	return config.MaxBlockCpuUsage
}

func (rl *ResourceLimitsManager) InitializeAccount(account name.AccountName) error {
	if err := rl.session.CreateResourceLimits(&resource.ResourceLimits{
		Owner:     account,
		NetWeight: -1,
		CpuWeight: -1,
		RamBytes:  -1,
	}); err != nil {
		return err
	}

	if err := rl.session.CreateResourceUsage(&resource.ResourceUsage{
		Owner: account,
	}); err != nil {
		return err
	}

	return nil
}

func (rl *ResourceLimitsManager) AddPendingRamUsage(account name.AccountName, ramDelta int64) error {
	if ramDelta == 0 {
		return nil
	}

	usage, err := rl.session.FindResourceUsageByOwner(account)

	if err != nil {
		return err
	}

	if ramDelta > 0 && math.MaxUint64-usage.RamUsage < uint64(ramDelta) {
		return fmt.Errorf("ram usage delta would overflow UINT64_MAX")
	} else if ramDelta < 0 && usage.RamUsage < uint64(-ramDelta) {
		return fmt.Errorf("ram usage delta would underflow UINT64_MAX")
	}

	// TODO: Add checks
	return rl.session.ModifyResourceUsage(usage, func() {
		usage.RamUsage += uint64(ramDelta)
	})
}

func (rl *ResourceLimitsManager) VerifyAccountRamUsage(account name.AccountName) error {
	var ramBytes, netWeight, cpuWeight *int64

	if err := rl.GetAccountLimits(account, ramBytes, netWeight, cpuWeight); err != nil {
		return err
	}

	usage, err := rl.session.FindResourceUsageByOwner(account)

	if err != nil {
		return err
	}

	if *ramBytes >= 0 {
		if usage.RamUsage > uint64(*ramBytes) {
			return fmt.Errorf("account %s has insufficient ram; needs %v bytes has %v bytes", account, usage.RamUsage, *ramBytes)
		}
	}

	return nil
}

func (rl *ResourceLimitsManager) GetAccountLimits(account name.AccountName, ramBytes *int64, netWeight *int64, cpuWeight *int64) error {
	pendingBuo, _ := rl.session.FindResourceLimitsByOwner(true, account)

	if pendingBuo != nil {
		*ramBytes = pendingBuo.RamBytes
		*netWeight = pendingBuo.NetWeight
		*cpuWeight = pendingBuo.CpuWeight
	} else {
		buo, err := rl.session.FindResourceLimitsByOwner(true, account)

		if err != nil {
			return err
		}

		*ramBytes = buo.RamBytes
		*netWeight = buo.NetWeight
		*cpuWeight = buo.CpuWeight
	}

	return nil
}
