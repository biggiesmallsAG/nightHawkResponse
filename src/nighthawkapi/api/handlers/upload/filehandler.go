package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	api "nighthawkapi/api/core"

	"os"
)

func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder must be set to check UserID/cookie/auth before placing job in upload queue; DE
	err := r.ParseMultipartForm(api.UPLOAD_MEM)
	if err != nil {
		api.LogError(api.DEBUG, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := r.MultipartForm
	case_name := m.Value["case_name"]

	for _, fheaders := range m.File {
		for _, hdr := range fheaders {
			file, err := hdr.Open()

			defer file.Close()

			if err != nil {
				api.LogError(api.DEBUG, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			dst, err := os.Create(fmt.Sprintf("%s/%s", api.MEDIA_DIR, hdr.Filename))
			defer dst.Close()

			if err != nil {
				api.LogError(api.DEBUG, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, file); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
	response := api.Success{
		Reason:   fmt.Sprintf("Success: File(s) Uploaded.. Ref Case: %s, check job queue for results.", case_name[0]),
		Response: http.StatusText(200),
	}

	ret, err := json.MarshalIndent(&response, "", "    ")

	if err != nil {
		api.LogError(api.DEBUG, err)
	}
	api.LogDebug(api.DEBUG, "[+] POST /upload HTTP 200, files uploaded. Generating dispatch for queue.")

	JobDispatch(m, case_name)
	fmt.Fprintln(w, string(ret))
}
