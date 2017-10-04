package analyze

import (
	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	"nighthawk/stack"
)



func ProcessIsVerified(process nhs.ProcessItem) (bool, string) {
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {
		// Checking Process properties in stack database
		// Arguments passed are
		// audit: w32processes-memory
		// name: ProcessItem.Name
		// path: ProcessItem.Path
		// regpath: ProcessItem.Arguments || ""
		// additional_info: ""
		if stack.IsCommonStackItem("w32processes-memory", process.Name, process.Path, "", "") {
			return true, "Verified by stacking processes"
		}
	}

	// default return 
	return false, ""
}