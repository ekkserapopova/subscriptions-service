package reader

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReadResponseData(r *http.Request, request interface{}) error {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	defer r.Body.Close()
	if err := json.Unmarshal(data, &request); err != nil {
		return err
	}

	return nil
}
