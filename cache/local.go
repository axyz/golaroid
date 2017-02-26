package cache

import (
	"path/filepath"

	"github.com/axyz/golaroid"
	"github.com/axyz/golaroid/reader"
	"github.com/axyz/golaroid/writer"
)

func get(variant string, path string, opts map[string]interface{}) ([]byte, error) {
	return reader.Local(variant, filepath.Join(variant, path), opts)
}

var Local = golaroid.Cache{
	Get: get,
	Set: writer.Local,
}
