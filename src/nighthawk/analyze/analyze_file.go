package analyze

import (
	nhconfig "nighthawk/config"	
	nhs "nighthawk/nhstruct"
	"nighthawk/hashset"
)


func FileIsBlacklisted(file *nhs.FileItem) (bool) {
	//// Check BlackList
	fi := nhs.BlacklistItem{
		AuditType: "w32files",
		Name: file.FileName,
		Path: file.FilePath,
		Md5sum: file.Md5sum,
	}
	
	return QueryBlacklistInformation(&fi)
}

func RawFileIsBlacklisted(file *nhs.RawFileItem) (bool) {
	//// Check BlackList
	fi := nhs.BlacklistItem{
		AuditType: "w32files",
		Name: file.FileName,
		Path: file.FilePath,
		Md5sum: file.Md5sum,
	}
	
	return QueryBlacklistInformation(&fi)
}

func FileIsWhitelisted(file *nhs.FileItem) bool {
	//// Check Whitelist
	fi := nhs.WhitelistItem{
		AuditType: "w32files",
		Name: file.FileName,
		Path: file.FilePath,
		Md5sum: file.Md5sum,
	}

	return QueryWhitelistInformation(&fi)
}

func RawFileIsWhitelisted(file *nhs.RawFileItem) bool {
	//// Check Whitelist
	fi := nhs.WhitelistItem{
		AuditType: "w32files",
		Name: file.FileName,
		Path: file.FilePath,
		Md5sum: file.Md5sum,
	}

	return QueryWhitelistInformation(&fi)
}

func FileIsVerified(file *nhs.FileItem) (bool, string) {
	//// Check file Known signer
	// If binary is signed and verified and reputed signer
	// consider it as good binary
	if IsKnownCertIssuer(file.PeInfo.DigitalSignature.CertificateIssuer) {
		return true, "verified by known digital signature signer"
	}

	// If HashSet checking is not enabled in nighthawk.json configuration
	// or HastSet database is not available in Elasticserach 
	// then skip checking HashSet
	if nhconfig.HashSetEnabled() && nhconfig.HashSetAvailable() {
		// If a match is found in NSRL HashSet return true 
		hashset.LoadHashSetConfig()
		if hashset.SearchWhitelistHash("md5", file.Md5sum, 10, false) {
			return true, "verified by NSRL HashSet"
		} 
	}

	// Default return
	return false, ""
}


func RawFileIsVerified(file *nhs.RawFileItem) (bool, string) {
	//// Check file Known signer
	// If binary is signed and verified and reputed signer
	// consider it as good binary
	if IsKnownCertIssuer(file.PeInfo.DigitalSignature.CertificateIssuer) {
		return true, "verified by known digital signature signer"
	}



	// If HashSet checking is not enabled in nighthawk.json configuration
	// or HastSet database is not available in Elasticserach 
	// then skip checking HashSet
	if nhconfig.HashSetEnabled() && nhconfig.HashSetAvailable() {
		// If a match is found in NSRL HashSet return true 
		hashset.LoadHashSetConfig()
		if hashset.SearchWhitelistHash("md5", file.Md5sum, 10, false) {
			return true, "verified by NSRL HashSet"
		} 
	}

	// Default return
	return false, ""
}

