# golaroid
media transformation service

## Usage
Basic usage with common included modules
```go
package main

import (
	"runtime"

	"github.com/axyz/golaroid"
	"github.com/axyz/golaroid/cache"
	"github.com/axyz/golaroid/reader"
	"github.com/axyz/golaroid/transform"
	"github.com/axyz/golaroid/writer"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	options := golaroid.Options{
		Port: 8080,
		Readers: golaroid.Readers{
			"local": reader.Local,
		},
		Transformers: golaroid.Transformers{
			"resize": transform.Resize,
			"strip":  transform.Strip,
			"jpeg":   transform.Jpeg,
		},
		Writers: golaroid.Writers{
			"local": writer.Local,
		},
		Caches: golaroid.Caches{
			"local": cache.Local,
		},
	}

	golaroid.Run(options)
}
```

example config:
```yaml
entities:
- Name: media
  Route: "/media/"

  InputReader:
    name: local
    options:
      folder: "./media"

  OutputWriters:
  - name: local
    options:
      folder: "./output"

  CachingLayer:
    name: local
    options:
      folder: "./cache"

  Variants:
    thumb:
      transformations:
      - name: resize
        options:
          width: 200
          height: 600
          filter: lanczos
          multistep: true
      - name: strip
      - name: jpeg
        options:
          quality: 85
```

## Modules
golaroid is extremely modular and custom `Readers`, `Writers`, 
`Caches` and `Transformers` for each entity can be defined and plugged as
required.

Some basic modules are provided with the golaroid package itself:
- `golaroid.reader`
  - `Local`: simple local file reader 
- `golaroid.writer`
  - `Local`: simple local file writer
- `golaroid.cache`
  - `Local`: simple file based cache
- `golaroid.transform`
  - `Jpeg`: jpeg encoding using imageMagick
  - `Resize`: resizing using imageMagick with multi-step support
  - `Strip`: metadata removal for size reduction
  
Writing your own modules is also possible and simple as you only have to define
functions operating on binary data.

### Readers
A `Reader` should have type
```go
func(variant string, path string, opts map[string]interface{}) ([]byte, error)
```
Given the variant name, a path for the media item and arbitrary options should return
the binary data for the media item and an error object.

### Writers
A `Writer` should have type
```go
func(variant string, path string, data []byte, opts map[string]interface{}) error
```
Given the variant name, a path for the media item, the media item itself as
binary data and arbitrary options should write somewhere the data by side effect
and return an error object.

### Caches
a `Cache` should be a `golaroid.Cache` struct that defines `Get` and a `Set`
methods.
```go
type Cache struct {
	Get func(variant string, path string, opts map[string]interface{}) ([]byte, error)
	Set func(variant string, path string, data []byte, opts map[string]interface{}) error
}
```
`Get` and `Set` methods should allow to save and retrieve binary data for the
media items given the variant, the path and arbitrary options.
