/*
 * nighthawkresponse - api server
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	nhr "nighthawkresponse"
	api "nighthawkresponse/api/core"
	"nighthawkresponse/api/handlers/analyzer"
	"nighthawkresponse/api/handlers/audit"
	"nighthawkresponse/api/handlers/auth"
	"nighthawkresponse/api/handlers/config"
	"nighthawkresponse/api/handlers/delete"
	"nighthawkresponse/api/handlers/search"
	"nighthawkresponse/api/handlers/stacking"
	"nighthawkresponse/api/handlers/upload"
	"nighthawkresponse/api/handlers/watcher"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

var logFile *os.File
var err error

func init() {
	logDir := filepath.Join(api.STATEDIR, "log")
	logFileName := filepath.Join(logDir, "access.log")
	logFile, err = os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// ServerConfig API Server runtime options
type ServerConfig struct {
	Server  string
	Port    int
	Version bool
	Debug   bool
	Help    bool
}

func fUsage() {
	fmt.Printf("\tnightHawkAPI v%s, by Team nightHawk (Daniel Eden & Roshan Maskey).\n", api.VERSION)
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(0)
}

func main() {
	flag.Usage = fUsage
	var svrconf ServerConfig

	flag.StringVar(&svrconf.Server, "server", "localhost", "Bind server to address")
	flag.IntVar(&svrconf.Port, "port", 8080, "Bind server to port")
	flag.BoolVar(&svrconf.Debug, "debug", false, "Turn on console level debbuging")
	flag.BoolVar(&svrconf.Help, "help", false, "Display use flags")
	flag.BoolVar(&svrconf.Version, "version", false, "Show version information")

	flag.Parse()

	if svrconf.Help {
		fUsage()
	}

	if svrconf.Version {
		nhr.ShowVersion("nighthawk response - API Server")
		os.Exit(0)
	}

	if svrconf.Debug {
		api.DEBUG = true
	}

	if svrconf.Server == "" {
		svrconf.Server = "localhost"
	}

	if svrconf.Port < 10 {
		svrconf.Port = 8080
	}

	router := mux.NewRouter()

	//----------------------------------------------------------------
	// Authentication API
	//----------------------------------------------------------------
	router.HandleFunc("/api/v1/auth/login", SessionHandler(auth.Login)).Methods("POST")
	router.HandleFunc("/api/v1/auth/password/change", SessionHandler(auth.ChangePassword)).Methods("POST")
	router.HandleFunc("/api/v1/auth/logout", SessionHandler(auth.Logout)).Methods("POST")

	//----------------------------------------------------------------
	// Admin APIs
	//----------------------------------------------------------------
	router.HandleFunc("/api/v1/admin/user/create", SessionHandler(auth.CreateNewUser)).Methods("POST")
	router.HandleFunc("/api/v1/admin/password/set", SessionHandler(auth.SetPassword)).Methods("POST")
	router.HandleFunc("/api/v1/admin/user/delete", SessionHandler(auth.DeleteUser)).Methods("POST")

	//----------------------------------------------------------------
	// Authenticated User APIs
	//----------------------------------------------------------------
	//// Platform specific APIs
	router.HandleFunc("/api/v1/config", SessionHandler(config.ReturnSystemConfig)).Methods("GET")
	router.HandleFunc("/api/v1/config", SessionHandler(config.UpdateSystemConfig)).Methods("POST")
	router.HandleFunc("/api/v1/platformstats", SessionHandler(config.ReturnPlatformStats)).Methods("GET")

	router.HandleFunc("/api/v1/system/config", SessionHandler(config.ReturnSystemConfig)).Methods("GET")
	router.HandleFunc("/api/v1/system/config", SessionHandler(config.UpdateSystemConfig)).Methods("POST")
	router.HandleFunc("/api/v1/system/stats", SessionHandler(config.ReturnPlatformStats)).Methods("GET")

	router.HandleFunc("/api/v1/upload", SessionHandler(upload.UploadFileHandler)).Methods("POST")
	router.HandleFunc("/api/v1/list/completedjobs", SessionHandler(upload.ListCompletedJobs)).Methods("GET") // deprecate
	router.HandleFunc("/api/v1/list/jobs/completed", SessionHandler(upload.ListCompletedJobs)).Methods("GET")
	router.HandleFunc("/api/v1/list/uploadjobs", SessionHandler(upload.SubscribeJobs)).Methods("GET") // deprecate
	router.HandleFunc("/api/v1/list/jobs/uploaded", SessionHandler(upload.SubscribeJobs)).Methods("GET")

	//// Basic Case/Endpoint/Audit management
	router.HandleFunc("/api/v1/list/cases", SessionHandler(audit.GetCaseList)).Methods("GET")
	router.HandleFunc("/api/v1/list/endpoints", SessionHandler(audit.GetEndpointList)).Methods("GET")
	router.HandleFunc("/api/v1/list/audittypes", SessionHandler(audit.GetAuditTypeList)).Methods("GET")
	router.HandleFunc("/api/v1/list/audits", SessionHandler(audit.GetAuditTypeList)).Methods("GET")

	router.HandleFunc("/api/v1/delete/case/{case}", SessionHandler(delete.DeleteCase)).Methods("GET")
	router.HandleFunc("/api/v1/delete/case/all", SessionHandler(delete.DeleteCase)).Methods("GET")
	router.HandleFunc("/api/v1/delete/endpoint/{endpoint}", SessionHandler(delete.DeleteEndpoint)).Methods("GET")
	router.HandleFunc("/api/v1/delete/endpont/all", SessionHandler(delete.DeleteEndpoint)).Methods("GET")
	router.HandleFunc("/api/v1/delete/{case}/{endpoint}", SessionHandler(delete.DeleteCaseEndpoint)).Methods("GET")
	// flexible delete option using POST
	router.HandleFunc("/api/v1/delete/case", SessionHandler(delete.DeleteCase)).Methods("POST")
	router.HandleFunc("/api/v1/delete/endpoint", SessionHandler(delete.DeleteEndpoint)).Methods("POST")

	//// View Case/Endpoint/Audits data
	router.HandleFunc("/api/v1/show/cases", SessionHandler(audit.GetCaseList)).Methods("GET")
	router.HandleFunc("/api/v1/show/endpoints", SessionHandler(audit.GetEndpointList)).Methods("GET")
	router.HandleFunc("/api/v1/show/audits", SessionHandler(audit.GetAuditTypeList)).Methods("GET")
	router.HandleFunc("/api/v1/show/doc/{doc_id}/{endpoint}", SessionHandler(audit.GetDocById)).Methods("GET")
	router.HandleFunc("/api/v1/show/{case}", SessionHandler(audit.GetEndpointByCase)).Methods("GET") // deprecate this api
	router.HandleFunc("/api/v1/show/{case}/endpoints", SessionHandler(audit.GetEndpointByCase)).Methods("GET")
	router.HandleFunc("/api/v1/show/{case}/{endpoint}", SessionHandler(audit.GetCasedateByEndpoint)).Methods("GET") // deprecate this
	router.HandleFunc("/api/v1/show/{case}/{endpoint}/dates", SessionHandler(audit.GetCasedateByEndpoint)).Methods("GET")
	router.HandleFunc("/api/v1/show/{case}/{endpoint}/{case_date}", SessionHandler(audit.GetAuditTypeByEndpointAndCase)).Methods("GET") // deprecate this
	router.HandleFunc("/api/v1/show/{case}/{endpoint}/{case_date}/audits", SessionHandler(audit.GetAuditTypeByEndpointAndCase)).Methods("GET")
	router.HandleFunc("/api/v1/show/{case}/{endpoint}/{case_date}/{audittype}", SessionHandler(audit.GetAuditDataByAuditGenerator)).Methods("GET", "POST")

	//// Stacking features
	router.HandleFunc("/api/v1/stacking/context", SessionHandler(stacking.GetStackContext)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/context/endpoint", SessionHandler(stacking.GetStackContextByEndpoint)).Methods("POST")

	router.HandleFunc("/api/v1/stacking/service", SessionHandler(stacking.StackServices)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/prefetch", SessionHandler(stacking.StackPrefetch)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/task", SessionHandler(stacking.StackTasks)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/persistence", SessionHandler(stacking.StackPersistence)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/locallistenport", SessionHandler(stacking.StackLocalListenPort)).Methods("POST") // deprecate
	router.HandleFunc("/api/v1/stacking/ports/listening", SessionHandler(stacking.StackLocalListenPort)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/runkey", SessionHandler(stacking.StackRunKey)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/dns/a", SessionHandler(stacking.StackDnsARequest)).Methods("POST")
	router.HandleFunc("/api/v1/stacking/url/domain", SessionHandler(stacking.StackDnsARequest)).Methods("POST")

	//// Search Global/Timeline
	router.HandleFunc("/api/v1/search", SessionHandler(search.GetGlobalSearch)).Methods("POST")
	router.HandleFunc("/api/v1/search/timeline", SessionHandler(search.GetTimelineSearch)).Methods("POST")

	//// Endpoint Diffing
	router.HandleFunc("/api/v1/diff/{endpoint}", SessionHandler(stacking.TimelineEndpointDiff)).Methods("GET") // deprecate
	router.HandleFunc("/api/v1/diff", SessionHandler(stacking.TimelineEndpointDiff)).Methods("POST")           // deprecate
	router.HandleFunc("/api/v1/stacking/endpoint/{endpoint}", SessionHandler(stacking.TimelineEndpointDiff)).Methods("GET")
	router.HandleFunc("/api/v1/stacking/endpoint", SessionHandler(stacking.TimelineEndpointDiff)).Methods("POST")

	//// Watcher
	router.HandleFunc("/api/v1/watcher/generate/rule", SessionHandler(watcher.GenerateWatcherRule)).Methods("POST")
	router.HandleFunc("/api/v1/watcher/result/all", SessionHandler(watcher.GetWatcherResults)).Methods("GET")
	router.HandleFunc("/api/v1/watcher/result/{id}", SessionHandler(watcher.GetWatcherResultById)).Methods("GET")

	//// Blacklist/Whitelist/Stack-Whitelist
	router.HandleFunc("/api/v1/analyze/add/blacklist", SessionHandler(analyzer.AddBlacklistInformation)).Methods("POST")
	router.HandleFunc("/api/v1/analyze/add/whitelist", SessionHandler(analyzer.AddWhitelistInformation)).Methods("POST")
	router.HandleFunc("/api/v1/analyze/add/stack", SessionHandler(analyzer.AddStackInformation)).Methods("POST")
	router.HandleFunc("/api/v1/analyze/show/{analyzer_type}", SessionHandler(analyzer.ShowAnalyzerItemByType)).Methods("GET")
	router.HandleFunc("/api/v1/analyze/delete/{analyzer_type}/{analyzer_id}", SessionHandler(analyzer.DeleteAnalyzerItemByID)).Methods("GET")
	router.HandleFunc("/api/v1/analyze/delete/{analyzer_type}", SessionHandler(analyzer.DeleteAnalyzerItemByQuery)).Methods("POST")

	//// Document Comments
	router.HandleFunc("/api/v1/comment/case-and-computer", SessionHandler(audit.GetCommentCaseAndComputer)).Methods("GET")
	router.HandleFunc("/api/v1/comment/add/{case}/{endpoint}/{audit}/{doc_id}", SessionHandler(audit.AddComment)).Methods("POST")
	router.HandleFunc("/api/v1/comment/show", SessionHandler(audit.GetComment)).Methods("GET", "POST")
	router.HandleFunc("/api/v1/comment/show/{case}", SessionHandler(audit.GetComment)).Methods("GET")
	router.HandleFunc("/api/v1/comment/show/{case}/{endpoint}", SessionHandler(audit.GetComment)).Methods("GET")
	router.HandleFunc("/api/v1/comment/show/{case}/{endpoint}/{audit}", SessionHandler(audit.GetComment)).Methods("GET")
	router.HandleFunc("/api/v1/comment/show/{case}/{endpoint}/{audit}/{doc_id}", SessionHandler(audit.GetComment)).Methods("GET")

	//// Document Tags
	router.HandleFunc("/api/v1/tag/case-and-computer", SessionHandler(audit.GetTagCaseAndComputer)).Methods("GET")
	router.HandleFunc("/api/v1/tag/add/{case}/{endpoint}/{audit}/{doc_id}", SessionHandler(audit.AddTag)).Methods("POST")
	router.HandleFunc("/api/v1/tag/show", SessionHandler(audit.GetTagData)).Methods("GET", "POST")
	router.HandleFunc("/api/v1/tag/show/{case}", SessionHandler(audit.GetTagData)).Methods("GET")
	router.HandleFunc("/api/v1/tag/show/{case}/{endpoint}", SessionHandler(audit.GetTagData)).Methods("GET")
	router.HandleFunc("/api/v1/tag/show/{case}/{endpoint}/{audit}", SessionHandler(audit.GetTagData)).Methods("GET")
	router.HandleFunc("/api/v1/tag/show/{case}/{endpoint}/{audit}/{doc_id}", SessionHandler(audit.GetTagData)).Methods("GET")

	//// Start mux router
	http.Handle("/", router)
	logger := handlers.LoggingHandler(logFile, router)
	api.LogDebug(api.DEBUG, fmt.Sprintf("[-] Serving on %s", fmt.Sprintf("%s:%d", svrconf.Server, svrconf.Port)))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", svrconf.Server, svrconf.Port), logger))
}

// SessionHandler handles authenticationa and admin session
func SessionHandler(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.RequestURI, "/api/v1/auth/login") {
			// do nothing
		} else if strings.HasPrefix(r.RequestURI, "/api/v1/admin") {
			isadmin, message := auth.IsAuthenticatedAdminSession(w, r)
			if !isadmin {
				api.HttpResponseReturn(w, r, "failed", message, nil)
				return
			}
		} else {
			isauth, message := auth.IsAuthenticatedSession(w, r)
			if !isauth {
				api.HttpResponseReturn(w, r, "failed", message, nil)
				return
			}
		}
		f(w, r)
	}
}
