package analyze

import (
	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
)

func TaskIsBlacklisted(task *nhs.TaskItem) bool {
	ti := nhs.BlacklistItem{
		AuditType: "w32tasks",
		Name: task.Name,
		Path: task.Path,
		TaskCreator: task.Creator,
	}

	return QueryBlacklistInformation(&ti)
}


func TaskIsWhitelisted(task *nhs.TaskItem) bool {
	ti := nhs.WhitelistItem{
		AuditType: "w32tasks",
		Name: task.Name,
		Path: task.Path,
		TaskCreator: task.Creator,
	}

	return QueryWhitelistInformation(&ti)
}


func TaskIsVerified(task *nhs.TaskItem) (bool, string) {
	//// Check task Stacking
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {

		si := nhs.StackItem{
			AuditType:"w32tasks",
			Name: task.Name,
			Path: task.Path,
			TaskCreator: task.Creator,
		}

		if IsCommonStackItem(&si) {
			return true, "Verified by stacking tasks"
		}
	}

	// default return
	return false, ""
}