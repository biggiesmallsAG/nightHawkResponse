package stack

import (
	"database/sql"
	//"fmt"
	//"os"

	//nhc "nighthawk/common"
	nhconfig "nighthawk/config"
	nhlog "nighthawk/log"

	//_ "github.com/mattn/go-sqlite3"
)

var CAList []string
var db *sql.DB = nil

func init() {

	// If check_stack: false in nighthawk.json configuration
	// or nighthawk.db is not available then ignore
	// initializing and process sqlite3 database
	if !nhconfig.CheckStack() || !nhconfig.StackDbEnabled() || !nhconfig.StackDbAvailable() {
		return
	}
	//db = openStackDatabase("/opt/nighthawk/etc/nighthawk.db")
	/*
	db = nhconfig.StackDbObject()
	if db == nil {
		fmt.Println("Stack::init - db is nil")
	}
	*/
	CAList = dbGetCaList()
}

/*
func openStackDatabase(dbfile string) *sql.DB {
	sdb, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		nhlog.LogMessage("openStackDatabase", "ERROR", fmt.Sprintf("Failed to open stackdb %s. %s", dbfile, err.Error()))
		return nil
	}
	return sdb
}
*/

func dbGetCaList() []string {
	var calist []string

	/*
	sqlQuery := fmt.Sprintf("SELECT CertIssuer FROM KnownIssuer")
	rows, err := db.Query(sqlQuery)
	if err != nil {
		nhlog.LogMessage("dbGetCaList", "ERROR", err.Error())
		os.Exit(nhc.ERROR_STACK_CERT_QUERY)
	}
	defer rows.Close()

	var ci string
	for rows.Next() {
		rows.Scan(&ci)
		calist = append(calist, ci)
	}
	*/
	return calist
}

func IsKnownCertIssuer(certissuer string) bool {
	if certissuer == "" {
		return false
	}

	for _, ca := range CAList {
		if ca == certissuer {
			return true
		}
	}

	return false
}

func IsCommonStackItem(audit string, name string, path string, regpath string, additional_info string) bool {

	// Return false if all the parameters are null
	if name == "" && path == "" && regpath == "" && additional_info == "" {
		nhlog.LogMessage("IsCommonStackItem", "DEBUG", "Empty or NULL parameters")
		return false
	}

	/*
	sqlQuery := "SELECT Id FROM Stack WHERE"
	if audit == "" {
		nhlog.LogMessage("IsCommonStackItem", "ERROR", "AuditType value missing. This is mandatory field")
		os.Exit(nhc.ERROR_STACK_NO_AUDITTYPE)
	}

	sqlQuery = fmt.Sprintf("%s Audit='%s'", sqlQuery, audit)

	if name != "" {
		sqlQuery = fmt.Sprintf("%s AND Name='%s'", sqlQuery, name)
	}

	if path != "" {
		sqlQuery = fmt.Sprintf("%s AND Path='%s'", sqlQuery, path)
	}

	if regpath != "" {
		sqlQuery = fmt.Sprintf("%s AND RegPath='%s'", sqlQuery, regpath)
	}

	if additional_info != "" {
		sqlQuery = fmt.Sprintf("%s AND AdditionalInfo='%s'", sqlQuery, additional_info)
	}

	nhlog.LogMessage("IsCommonStackItem", "DEBUG", fmt.Sprintf("Query: %s", sqlQuery))

	rows, err := db.Query(sqlQuery)
	if err != nil {
		nhlog.LogMessage("IsCommonStackItem", "ERROR", fmt.Sprintf("Failed to query Stack table. %s", err.Error()))
		//os.Exit(nhc.ERROR_STACK_ITEM_QUERY)
		rows.Close()
		return false
	}

	var Id int
	for rows.Next() {
		rows.Scan(&Id)
		if Id > 0 {
			nhlog.LogMessage("IsCommonStackItem", "DEBUG", fmt.Sprintf("%s: OK", sqlQuery))
			return true
		}
	}
	rows.Close()
	*/
	return false
}
