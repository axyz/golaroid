package golaroid

import (
	"fmt"
	"net/http"
	"strings"
)

func readMedia(
	o Options,
	variantName string,
	path string,
	r inputReader) ([]byte, error) {

	fmt.Println("TODO: get image from path: " + path)

	m, err := o.Readers[r.Name](variantName, path, r.Options)
	return m, err
}

func transformMediaItem(
	o Options,
	data []byte,
	ts []transformation) ([]byte, error) {

	result := data

	// serially apply all the transformers to the original binary data
	for _, t := range ts {
		fmt.Println("TODO: apply transformation: " + t.Name)
		d, err := o.Transformers[t.Name](result, t.Options)
		result = d
		if err != nil {
			return result, err
		}
	}

	fmt.Println("TODO: return transformed result")
	return result, nil
}

func writeMedia(
	o Options,
	variantName string,
	data []byte,
	path string,
	ws outputWriters) error {

	// write data using all the defined writers
	// TODO: writes should run in parallel
	for _, w := range ws {
		err := o.Writers[w.Name](variantName, path, data, w.Options)
		if err != nil {
			return err
		}
		fmt.Println("writing data to " + w.Name)
	}

	return nil
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "{ \"ok\": true }")
	default:
		w.WriteHeader(405)
	}
}

func createHandler(o Options, c *EntityConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			p := r.URL.Path[len(c.Route):]
			_str := strings.SplitN(p, "/", 2)
			imgPath := _str[1]
			variantName := _str[0]
			variant := c.Variants[variantName]

			transformations := variant.Transformations
			reader := c.InputReader
			writers := c.OutputWriters
			cl := c.CachingLayer
			if cl.Name != "" {
				cache := o.Caches[cl.Name]
				data, err := cache.Get(variantName, imgPath, cl.Options)
				if err == nil {
					fmt.Println("serving from cache")
					w.Header().Set("Content-Type", "image/jpeg")
					w.Write(data)
					return
				}
			}

			m, err := readMedia(o, variantName, imgPath, reader)
			if err != nil {
				fmt.Println(err.Error())
			}

			tm, err := transformMediaItem(o, m, transformations)
			if err != nil {
				fmt.Println(err.Error())
			}

			writeMedia(o, variantName, tm, imgPath, writers)

			if cl.Name != "" {
				cache := o.Caches[cl.Name]
				fmt.Println("saving to cache")
				cache.Set(variantName, imgPath, tm, cl.Options)
			}

			w.Header().Set("Content-Type", "image/jpeg")
			w.Write(tm)
		default:
			// method not allowed
			w.WriteHeader(405)
		}
	}
}
