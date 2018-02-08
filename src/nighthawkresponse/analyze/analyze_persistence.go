package analyze

import (
	nhconfig "nighthawkresponse/config"
	nhs "nighthawkresponse/nhstruct"
)

func PersistenceIsBlacklisted(autorun *nhs.PersistenceItem) bool {
	//// Check process blacklist
	ai := nhs.BlacklistItem{
		AuditType:       "w32scripting-persistence",
		Name:            autorun.Registry.ValueName,
		Path:            autorun.Path,
		PersistenceType: autorun.PersistenceType,
		RegPath:         autorun.StackPath,
	}

	return QueryBlacklistInformation(&ai)
}

func PersistenceIsWhitelisted(autorun *nhs.PersistenceItem) bool {
	//// Check process blacklist
	ai := nhs.WhitelistItem{
		AuditType:       "w32scripting-persistence",
		Name:            autorun.Registry.ValueName,
		Path:            autorun.Path,
		PersistenceType: autorun.PersistenceType,
		RegPath:         autorun.StackPath,
	}

	return QueryWhitelistInformation(&ai)
}

func PersistenceIsVerified(autorun *nhs.PersistenceItem) (bool, string) {
	//// Check process stacking
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {

		ai := nhs.StackItem{
			AuditType:       "w32scripting-persistence",
			Name:            autorun.Registry.ValueName,
			Path:            autorun.Path,
			PersistenceType: autorun.PersistenceType,
			RegPath:         autorun.StackPath,
		}

		if IsCommonStackItem(&ai) {
			return true, "Verified by stacking persistence"
		}
	}

	// default return
	return false, ""
}
