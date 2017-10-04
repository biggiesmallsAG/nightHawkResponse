package analyze

import (
	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	"nighthawk/stack"
)


func TaskIsVerified(task nhs.TaskItem) (bool, string) {
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {
		// Arguments passed to IsCommonStackItem
		// audit: w32tasks
		// name: TaskItem.Name
		// path: TaskItem.Path
		// regpath: ""
		// additional_info: TaskItem.Creator 
		if stack.IsCommonStackItem("w32tasks", task.Name, task.Path,"",task.Creator) {
			return true, "Verified by stacking tasks"
		}
	}

	// default return
	return false, ""
}