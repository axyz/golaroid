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
	// a name to identify the entity e.g. media
	Name string

	// the base root for the entity e.g. /media/
	Route string

	// the reader layer to get media items
	InputReader inputReader

	// a list of writing layers to output generated media items
	OutputWriters []outputWriter

	// caching layer to be used
	CachingLayer cachingLayer

	// map of allowed variants for the current entity
	Variants map[string]variant
}

// arbitrary transformation over binary data
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

// golaroid server options
type Options struct {
	// listening port
	Port int

	// map of all available transformers
	Transformers

	// map of all available readers
	Readers

	// map of all available writers
	Writers

	// map of all available caches
	Caches
}

type entity struct {
	route   string
	handler func(http.ResponseWriter, *http.Request)
}

// define a route and a handler function to an entity given the EntityConfig
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

	// for each entityConfig create a new entity and push it to the
	// entities array
	for _, econf := range entitiesConfigs {
		e := new(entity).withConfig(o, &econf)
		entities = append(entities, e)
	}

	// register all the entities handlers
	for _, e := range entities {
		http.HandleFunc(e.route, e.handler)
	}

	// register healthcheck route
	http.HandleFunc(healthcheckRoute, healthcheck)

	// start the server
	http.ListenAndServe(":"+strconv.Itoa(o.Port), nil)
}
