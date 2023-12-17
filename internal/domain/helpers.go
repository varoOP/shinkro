package domain

import (
	"io"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

// func updateStart(ctx context.Context, s int) int {
// 	if s == 0 {
// 		return 1
// 	}
// 	return s
// }

func readYamlHTTP(resp *http.Response, mapping interface{}) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	err = yaml.Unmarshal(body, mapping)
	if err != nil {
		return err
	}

	return nil
}

func readYamlFile(f *os.File, mapping interface{}) error {
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, mapping)
	if err != nil {
		return err
	}

	return nil
}

func fileExists(path string) bool {
	_, err := os.Open(path)
	return err == nil
}
