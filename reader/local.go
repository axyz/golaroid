package reader

import (
	"io/ioutil"
	"path/filepath"
)

func Local(variant string, path string, opts map[string]interface{}) ([]byte, error) {
	folder := opts["folder"].(string)
	data, err := ioutil.ReadFile(filepath.Join(folder, path))
	return data, err
}
