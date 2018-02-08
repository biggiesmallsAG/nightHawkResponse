package elastic

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	nhc "nighthawkresponse/common"
	nhconfig "nighthawkresponse/config"
	nhlog "nighthawkresponse/log"
	nhs "nighthawkresponse/nhstruct"
)

var (
	client *http.Client
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func init() {
	client = createHttpClient()
}

func createHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func SplunkSessionKey(server string, port int, user string, password string) string {
	auth_url := fmt.Sprintf("https://%s:%d/services/auth/login?output_mode=json", server, port)
	auth_string := fmt.Sprintf("username=%s&password=%s", user, password)

	/*
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client := &http.Client{Transport: tr}
	*/
	req, _ := http.NewRequest("POST", auth_url, bytes.NewReader([]byte(auth_string)))
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Error authenticating to Splunk")
		fmt.Println("Status code: ", res.StatusCode)
		fmt.Println(res.Header)
		fmt.Println(res.Body)

		return ""
	}

	body, _ := ioutil.ReadAll(res.Body)

	var key map[string]string
	json.Unmarshal(body, &key)
	return key["sessionKey"]

}

func UploadToSplunk(computername string, audit string, data []nhs.RlRecord) {
	if len(data) == 0 {
		return
	}

	if nhconfig.ElasticBasicAuth() == "YWRtaW46YWRtaW4=" || nhconfig.ElasticBasicAuth() == "" {
		//fmt.Println("UploadToSplunk::ElasticBasicAuth()")
		authcode := SplunkSessionKey(nhconfig.ElasticHost(), nhconfig.ElasticPort(), nhconfig.ElasticUser(), nhconfig.ElasticPassword())
		nhconfig.SetAuthCode(authcode)
		if nhconfig.ElasticBasicAuth() == "" {
			nhlog.LogMessage("UploadToSplunk", "ERROR", "Failed to get Splunk session key")
			os.Exit(nhc.ERROR_SPLUNK_AUTHENTICATION)
		}
	}

	//fmt.Println("SessionKey: ", config.ElasticBasicAuth())

	upload_url := fmt.Sprintf("https://%s:%d/services/receivers/simple?source=nighthawkresponse&sourcetyp=redline&output_mode=json", nhconfig.ElasticHost(), nhconfig.ElasticPort())

	/*
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			MaxIdleConnsPerHost: MaxIdleConnection,
		}
		client = &http.Client{Transport: tr, Timeout: time.Duration(RequestTime)*time.Second}
	*/

	req, _ := http.NewRequest("POST", upload_url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Splunk "+nhconfig.ElasticBasicAuth())
	req.Header.Set("Connection", "keep-alive")

	for _, d := range data {
		jd, _ := json.Marshal(d)
		//fmt.Println(string(jd))
		//req, _ := http.NewRequest("POST", upload_url, bytes.NewReader(jd))
		//req.Header.Set("Accept", "application/json")
		//req.Header.Set("Authorization", "Splunk "+ config.ElasticBasicAuth())

		//body, ok := bytes.NewReader(jd).(io.ReadCloser)
		req.Body = ioutil.NopCloser(bytes.NewReader(jd))
		res, err := client.Do(req)
		if err != nil {
			//fmt.Println(err.Error())
			//fmt.Println(string(jd))
			nhlog.LogMessage("UploadToSplunk", "ERROR", err.Error())
		}

		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()

		if res.StatusCode != 200 {
			nhlog.LogMessage("UploadToSplunk", "INFO", fmt.Sprintf("ReqHeader: %s, StatusCode: %d, RespHeader: %s", req.Header, res.StatusCode, res.Header))
			os.Exit(nhc.ERROR_SPLUNK_UPLOAD)
		}
		//time.Sleep(1)
	}

}
