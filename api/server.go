package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud-ml/router"
	"cloud-ml/store"

	restful "github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

// APIServer ...
type APIServer struct {
	Config *APIServerOptions
}

const (
	// MongoDBHost ...
	MongoDBHost = "MONGODB_HOST"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// APIServerOptions contains all options(config) for api server
type APIServerOptions struct {
	MongoDBHost      string
	MongoGracePeriod time.Duration
	Port             int
	AddrTemplate     string
}

// NewAPIServerOptions returns a new APIServerOptions
func NewAPIServerOptions() *APIServerOptions {
	return &APIServerOptions{
		MongoDBHost:      "192.168.55.7",
		MongoGracePeriod: 30 * time.Second,
		Port:             7099,
		AddrTemplate:     "http://localhost:%v",
	}
}

// PrepareRun prepare for apiserver running
func (s *APIServer) PrepareRun() (*PreparedAPIServer, error) {
	closing := make(chan struct{})

	// init database
	_, err := store.Init(s.Config.MongoDBHost, s.Config.MongoGracePeriod, closing)
	if err != nil {
		return nil, err
	}

	return &PreparedAPIServer{s}, nil
}

// PreparedAPIServer is a prepared api server
type PreparedAPIServer struct {
	*APIServer
}

// Run start a api server
func (s *PreparedAPIServer) Run() error {
	dataStore := store.NewStore()
	defer dataStore.Close()

	// Initialize the V1 API.
	if err := router.InitRouters(dataStore); err != nil {
		return err
	}

	// start server
	server := &http.Server{Addr: fmt.Sprintf(":%d", s.Config.Port), Handler: restful.DefaultContainer}
	server.ListenAndServe()
	return nil
}

func main() {
	ops := NewAPIServerOptions()
	as := &APIServer{ops}
	ps, err := as.PrepareRun()
	if err != nil {
		log.Errorf("APIServer start failed,err msg:%s", err)
		return
	}
	log.Info("APIServer start succeed")
	ps.Run()
}
