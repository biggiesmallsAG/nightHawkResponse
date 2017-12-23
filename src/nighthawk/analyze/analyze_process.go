package analyze

import (
	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
)

func ProcessIsBlacklisted(process *nhs.ProcessItem) bool {
	//// Check process blacklist
	ai := nhs.BlacklistItem{
		AuditType: "w32processes",
		Name: process.Name,
		Path: process.Path,
		Arguments: process.Arguments,
	}

	return QueryBlacklistInformation(&ai)
}


func ProcessIsWhitelisted(process *nhs.ProcessItem) bool {
	//// Check process blacklist
	ai := nhs.WhitelistItem{
		AuditType: "w32processes",
		Name: process.Name,
		Path: process.Path,
		Arguments: process.Arguments,
	}

	return QueryWhitelistInformation(&ai)
}


func ProcessIsVerified(process *nhs.ProcessItem) (bool, string) {
	//// Check process stacking
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {

		var si nhs.StackItem 
		si.AuditType = "w32processes"
		si.Name = process.Name
		si.Path = process.Path

		if IsCommonStackItem(&si) {
			return true, "Verified by stacking processes"
		}
	}

	// default return 
	return false, ""
}