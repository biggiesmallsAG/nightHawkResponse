package analyze

import (
	nhconfig "nighthawkresponse/config"
	nhs "nighthawkresponse/nhstruct"
)

func ServiceIsBlacklisted(service *nhs.ServiceItem) bool {
	si := nhs.BlacklistItem{
		AuditType: "w32services",
		Name:      service.Name,
		Path:      service.Path,
		ServiceDescriptiveName: service.DescriptiveName,
	}

	return QueryBlacklistInformation(&si)
}

func ServiceIsWhitelisted(service *nhs.ServiceItem) bool {
	si := nhs.WhitelistItem{
		AuditType: "w32services",
		Name:      service.Name,
		Path:      service.Path,
		ServiceDescriptiveName: service.DescriptiveName,
	}

	return QueryWhitelistInformation(&si)
}

func ServiceIsVerified(service *nhs.ServiceItem) (bool, string) {

	//// Check service stacking
	// Check if the service details are registered to known commonly
	// configured service
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {
		si := nhs.StackItem{AuditType: "w32services",
			Name: service.Name,
			Path: service.Path,
			ServiceDescriptiveName: service.DescriptiveName,
		}
		if IsCommonStackItem(&si) {
			return true, "Verified by stacking services"
		}
	}
	// default
	return false, ""
}
