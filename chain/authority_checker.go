package chain

import (
	"github.com/MetalBlockchain/antelopevm/core/authority"
	"github.com/MetalBlockchain/antelopevm/crypto/ecc"
)

type PermissionCacheStatus int

const (
	BeingEvaluated PermissionCacheStatus = iota
	PermissionUnsatisfied
	PermissionSatisfied
)

type PermissionCacheType = map[authority.PermissionLevel]PermissionCacheStatus
type PermissionToAuthorityFunc func(*authority.PermissionLevel) (*authority.Authority, error)

type AuthorityChecker struct {
	PermissionToAuthority PermissionToAuthorityFunc
	ProvidedKeys          ecc.PublicKeySet
	ProvidedPermissions   []authority.PermissionLevel
	UsedKeys              ecc.PublicKeySet
	RecursionDepthLimit   uint16
}

func NewAuthorityChecker(permissionToAuthority PermissionToAuthorityFunc, keys ecc.PublicKeySet, providedPermissions []authority.PermissionLevel, recursionDepthLimit uint16) *AuthorityChecker {
	return &AuthorityChecker{
		PermissionToAuthority: permissionToAuthority,
		ProvidedKeys:          keys,
		ProvidedPermissions:   providedPermissions,
		UsedKeys:              ecc.NewPublicKeySet(keys.Size()),
		RecursionDepthLimit:   recursionDepthLimit,
	}
}

func (ac *AuthorityChecker) SatisfiedPermissionLevel(permission authority.PermissionLevel, cachedPerms *PermissionCacheType) bool {
	cachedPermissions := make(PermissionCacheType)

	if cachedPerms == nil {
		cachedPerms = ac.initializePermissionCache(&cachedPermissions)
	}

	visitor := NewWeightTallyVisitor(ac, cachedPerms, 0)

	return (visitor.Visit(authority.PermissionLevelWeight{Permission: permission, Weight: 1}) > 0)
}

func (ac *AuthorityChecker) SatisfiedAuthority(authority *authority.Authority, cachedPerms *PermissionCacheType) bool {
	cachedPermissions := make(PermissionCacheType)

	if cachedPerms == nil {
		cachedPerms = ac.initializePermissionCache(&cachedPermissions)
	}

	return ac.satisfiedAuthority(authority, cachedPerms, 0)
}

func (ac *AuthorityChecker) satisfiedAuthority(authority *authority.Authority, cachedPerms *PermissionCacheType, depth uint16) bool {
	permissions := make(MetaPermission, 0)

	for _, key := range authority.Keys {
		permissions = append(permissions, key)
	}

	for _, account := range authority.Accounts {
		permissions = append(permissions, account)
	}

	// Sort permissions by weight
	permissions.Sort()

	visitor := NewWeightTallyVisitor(ac, cachedPerms, depth)

	for _, permission := range permissions {
		if visitor.Visit(permission) >= authority.Threshold {
			return true
		}
	}

	return false
}

func (ac *AuthorityChecker) PermissionStatusInCache(permissions PermissionCacheType, level *authority.PermissionLevel) PermissionCacheStatus {
	value, ok := permissions[*level]

	if ok {
		return value
	}

	return 0
}

func (ac *AuthorityChecker) initializePermissionCache(cachedPermissions *PermissionCacheType) *PermissionCacheType {
	for _, p := range ac.ProvidedPermissions {
		(*cachedPermissions)[p] = PermissionSatisfied
	}

	return cachedPermissions
}

func (ac *AuthorityChecker) AllKeysUsed() bool {
	return ac.ProvidedKeys.Size() == ac.UsedKeys.Size()
}

func (ac *AuthorityChecker) GetUsedKeys() []ecc.PublicKey {
	return ac.UsedKeys.Slice()
}

func (ac *AuthorityChecker) GetUnusedKeys() []ecc.PublicKey {
	return ac.ProvidedKeys.Difference(ac.UsedKeys).Slice()
}

type WeightTallyVisitor struct {
	Checker           *AuthorityChecker
	CachedPermissions *PermissionCacheType
	TotalWeight       uint32
	RecursionDepth    uint16
}

func NewWeightTallyVisitor(checker *AuthorityChecker, cachedPermissions *PermissionCacheType, recursionDepth uint16) *WeightTallyVisitor {
	return &WeightTallyVisitor{
		Checker:           checker,
		CachedPermissions: cachedPermissions,
		RecursionDepth:    recursionDepth,
	}
}

func (w *WeightTallyVisitor) Visit(permission interface{}) uint32 {
	switch v := permission.(type) {
	case authority.KeyWeight:
		w.VisitKeyWeight(v)
		return w.TotalWeight
	case authority.PermissionLevelWeight:
		w.VisitPermissionLevelWeight(v)
		return w.TotalWeight
	default:
		return w.TotalWeight
	}
}

func (w *WeightTallyVisitor) VisitKeyWeight(permission authority.KeyWeight) uint32 {
	for _, key := range w.Checker.ProvidedKeys.Slice() {
		if key.Compare(permission.Key) {
			w.Checker.UsedKeys.Insert(key)
			w.TotalWeight += uint32(permission.Weight)
			break
		}
	}

	return w.TotalWeight
}

func (w *WeightTallyVisitor) VisitPermissionLevelWeight(permission authority.PermissionLevelWeight) uint32 {
	status := w.Checker.PermissionStatusInCache(*w.CachedPermissions, &permission.Permission)

	if status == BeingEvaluated {
		if w.RecursionDepth < w.Checker.RecursionDepthLimit {
			result := false
			auth, err := w.Checker.PermissionToAuthority(&permission.Permission)

			if err == nil && auth.Threshold > 0 {
				(*w.CachedPermissions)[permission.Permission] = BeingEvaluated
				result = w.Checker.satisfiedAuthority(auth, w.CachedPermissions, w.RecursionDepth+1)
			}

			if result {
				w.TotalWeight += uint32(permission.Weight)
				(*w.CachedPermissions)[permission.Permission] = PermissionSatisfied
			} else {
				(*w.CachedPermissions)[permission.Permission] = PermissionUnsatisfied
			}
		}
	} else if status == PermissionSatisfied {
		w.TotalWeight += uint32(permission.Weight)
	}

	return w.TotalWeight
}

type MetaPermission []interface{}

func (m MetaPermission) Len() int { return len(m) }

func (m MetaPermission) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MetaPermission) Less(i, j int) bool {
	iType, jType := 0, 0
	iWeight, jWeight := uint16(0), uint16(0)
	switch v := m[i].(type) {
	case authority.KeyWeight:
		iWeight = uint16(v.Weight)
		iType = 2
	case authority.PermissionLevelWeight:
		iWeight = uint16(v.Weight)
		iType = 3
	}
	switch v := m[j].(type) {
	case authority.KeyWeight:
		jWeight = uint16(v.Weight)
		iType = 2
	case authority.PermissionLevelWeight:
		jWeight = uint16(v.Weight)
		iType = 3
	}

	if iWeight < jWeight {
		return true
	} else if iWeight > jWeight {
		return false
	} else {
		if iType < jType {
			return true
		} else {
			return false
		}
	}
}

func (m MetaPermission) Sort() {
	for i := 0; i < m.Len()-1; i++ {
		for j := 0; j < m.Len()-1-i; j++ {
			if m.Less(j, j+1) {
				m.Swap(j, j+1)
			}
		}
	}
}
