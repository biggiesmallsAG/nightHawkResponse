package analyze

import (
	"fmt"
	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"

	elastic "gopkg.in/olivere/elastic.v5"
)

var (
	err    error
	query  elastic.Query
	client *elastic.Client
)

func init() {
	conf, err := nhconfig.LoadConfigFile(nhconfig.CONFIG_FILE)
	if err != nil {
		nhc.ConsoleMessage("Package analyze", "ERROR", err.Error(), true)
	}

	// Elasticsearch client initialization
	httpSchema := "http"
	if conf.Elastic.Ssl {
		httpSchema = "https"
	}

	client, err = elastic.NewClient(elastic.SetURL(fmt.Sprintf("%s://%s:%d", httpSchema, conf.Elastic.Host, conf.Elastic.Port)))
	if err != nil {
		nhc.ConsoleMessage("Package analyze", "ERROR", err.Error(), true)
		return
	}
}
