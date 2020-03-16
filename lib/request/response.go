package request

import "encoding/json"

type Response struct {
	Success  bool
	Message  string
	Data     interface{}
	Redirect string
}

func (r *Request) ResponseJSON(success bool, msg string, data interface{}, redirect string) error {
	if r.terminated {
		return nil
	}
	b, err := json.Marshal(Response{
		success, msg, data, redirect,
	})
	if err == nil {
		r.writer.Header().Set("Content-Type", "application/json")
		r.writer.Write(b)
		r.Terminate()
		return nil
	} else {
		return err
	}
}

func (r *Request) Response(success bool, msg string, data interface{}, redirect string) error {
	if r.terminated {
		return nil
	}
	if r.ContentType == REQ_JSON {
		data, err := json.Marshal(Response{
			success, msg, data, redirect,
		})
		if err == nil {
			r.writer.Header().Set("Content-Type", "application/json")
			r.writer.Write(data)
			r.Terminate()
			return nil
		} else {
			return err
		}
	} else {
		if success {
			if msg != "" {
				r.Success(msg)
			}
			r.Redirect(redirect)
		} else {
			if msg != "" {
				r.Error(msg)
			}
			redirect = r.request.URL.Path
			if r.request.URL.RawQuery != "" {
				redirect += "?" + r.request.URL.RawQuery
			}
			r.Redirect(redirect)
		}
	}
	return nil

}
