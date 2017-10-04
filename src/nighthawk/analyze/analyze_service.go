package analyze

import (
	nhconfig "nighthawk/config"
	nhs "nighthawk/nhstruct"
	"nighthawk/stack"
)

func ServiceIsVerified(service nhs.ServiceItem) (bool,string) {

	// Check if the service details are registered to known commonly
	// configured service
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {
		if stack.IsCommonStackItem("w32services", service.Name, service.Path, "", "") {
			return true, "Verified by stacking services"
		}
	}
	// default 
	return false,""
}