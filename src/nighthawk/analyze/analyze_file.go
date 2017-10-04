package analyze

import (
	"fmt"


	nhconfig "nighthawk/config"	
	nhs "nighthawk/nhstruct"
	nhlog "nighthawk/log"
	"nighthawk/hashset"
	"nighthawk/stack"
)



func FileIsVerified(file nhs.FileItem) (bool, string) {

	// If binary is signed and verified and reputed signer
	// consider it as good binary

	if stack.IsKnownCertIssuer(file.PeInfo.DigitalSignature.CertificateIssuer) {
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


	// Check if filename, filepath and md5 sum is in common persistence
	// database. 
	if nhconfig.StackDbEnabled() && nhconfig.StackDbAvailable() {
		
		nhlog.LogMessage("FileIsVerified", "DEBUG", fmt.Sprintf("Checking filestack %s", file.Path))
		

		if stack.IsCommonStackItem("w32scripting-persistence", "", file.Path,"","") {
			return true, "Verified by stacking persistence binaries"
		}
	}


	// Default return
	return false, ""
}


