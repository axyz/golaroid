package golaroid

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	healthcheckRoute = "/healthcheck"
)

type MediaItem struct {
	data        []byte
	name        string
	path        string
	contentType string
}

type action struct {
	Name    string
	Options map[string]interface{}
}

type transformation action
type inputReader action
type outputWriter action
type outputWriters []outputWriter
type cachingLayer action

type variant struct {
	Transformations []transformation
}

type EntityConfig struct {
	Name          string
	Route         string
	InputReader   inputReader
	OutputWriters []outputWriter
	CachingLayer  cachingLayer
	Variants      map[string]variant
}

type Transformer func(data []byte, opts map[string]interface{}) ([]byte, error)
type Transformers map[string]Transformer

type Reader func(variant string, path string, opts map[string]interface{}) ([]byte, error)
type Readers map[string]Reader

type Writer func(variant string, path string, data []byte, opts map[string]interface{}) error
type Writers map[string]Writer

type Cache struct {
	Get func(variant string, path string, opts map[string]interface{}) ([]byte, error)
	Set func(variant string, path string, data []byte, opts map[string]interface{}) error
}
type Caches map[string]Cache

type Options struct {
	Port int
	Transformers
	Readers
	Writers
	Caches
}

type entity struct {
	route   string
	handler func(http.ResponseWriter, *http.Request)
}

func (e entity) withConfig(o Options, c *EntityConfig) entity {
	e.route = c.Route
	e.handler = createHandler(o, c)
	return e
}

func Run(o Options) {
	viper.SetConfigName("config")
	viper.AddConfigPath(os.Getenv("GOLAROID_CONFIG_PATH"))
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("No config file found.")
	}
	viper.SetDefault("entities", []EntityConfig{})
	viper.WatchConfig()

	entities := []entity{}
	var entitiesConfigs []EntityConfig
	viper.UnmarshalKey("entities", &entitiesConfigs)

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		viper.UnmarshalKey("entities", &entitiesConfigs)
	})

	for _, econf := range entitiesConfigs {
		e := new(entity).withConfig(o, &econf)
		entities = append(entities, e)
	}

	for _, e := range entities {
		http.HandleFunc(e.route, e.handler)
	}

	http.HandleFunc(healthcheckRoute, healthcheck)

	http.ListenAndServe(":"+strconv.Itoa(o.Port), nil)
}
