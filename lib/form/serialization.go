package form

import (
	"encoding/json"
	"github.com/gorilla/schema"
	"github.com/iesreza/foundation/lib"
	"github.com/iesreza/foundation/lib/router"
	"github.com/iesreza/foundation/system"
	"strings"
)

func Unmarshal(req router.Request, output interface{}) error {
	//guess data type
	if req.Req().Header.Get("Content-Type") == "application/json" {
		return json.NewDecoder(req.Req().Body).Decode(output)
	}
	if req.Req().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		req.Req().ParseForm()
		return schema.NewDecoder().Decode(output, req.Req().Form)
	}
	if strings.HasPrefix(req.Req().Header.Get("Content-Type"), "multipart/form-data") {
		maxupload, _ := lib.ParseSize(system.GetConfig().App.MaxUploadSize)
		err := req.Req().ParseMultipartForm(int64(maxupload))
		if err != nil {
			return err
		}
		return schema.NewDecoder().Decode(output, req.Req().MultipartForm.Value)

	}
	return schema.NewDecoder().Decode(output, req.Req().Header)
}
