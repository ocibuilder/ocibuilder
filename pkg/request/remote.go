package request

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/beval/beval/pkg/apis/beval/v1alpha1"
)

func RequestRemote(url string, filepath string, auth v1alpha1.RemoteCreds) error {

	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if isRemoteAuth(auth) {
		if auth.Plain.Username != "" {
			req.SetBasicAuth(auth.Plain.Username, auth.Plain.Password)
		}

		if auth.Env.Username != "" {
			req.SetBasicAuth(os.Getenv(auth.Env.Username), os.Getenv(auth.Env.Password))
		}
	}

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http error received requesting remote overlay %d with response %s", res.StatusCode, res.Status)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			return
		}
	}()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}

	return nil
}

func isRemoteAuth(creds v1alpha1.RemoteCreds) bool {
	if creds.Plain.Username == "" && creds.Env.Username == "" {
		return false
	}
	return true
}
