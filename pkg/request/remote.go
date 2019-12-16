package request

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/docker/docker/api/types"
)

func RequestRemote(url string, filepath string, auth types.AuthConfig) error {

	cli := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if isAuthPresent(auth) {
		authString, err := generateAuthRegistryString(auth)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", authString)
	}

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http error received requesting remote overlay %d with response %s", res.StatusCode, res.Status)
	}

	defer func() {
		res.Body.Close()
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

func generateAuthRegistryString(auth types.AuthConfig) (string, error) {
	encodedJSON, err := json.Marshal(auth)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}

func isAuthPresent(auth types.AuthConfig) bool {

	if auth.Username != "" {
		return true
	}

	if auth.Auth != "" {
		return true
	}

	if auth.IdentityToken != "" {
		return true
	}

	return false
}
