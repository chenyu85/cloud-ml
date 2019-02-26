package store

import (
	"time"

	"cloud-ml/wait"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

const (
	defaultDBName         string = "cloud-ml"
	projectCollectionName string = "projects"
	socketTimeout                = time.Second * 5
	syncTimeout                  = time.Second * 5
	tickerDuration               = time.Second * 5
)

var (
	session *mgo.Session
	mclosed chan struct{}
)

// DataStore is the type for mongo db store.
type DataStore struct {
	s       *mgo.Session
	saltKey string

	// Collections
	projectCollection *mgo.Collection
}

// Init store mongo client session
func Init(host string, gracePeriod time.Duration, closing chan struct{}) (chan struct{}, error) {
	mclosed = make(chan struct{})
	var err error

	// dail mongo session
	// wait mongodb set up
	wait.Poll(time.Second, gracePeriod, func() (bool, error) {
		session, err = mgo.Dial(host)
		return err == nil, nil
	})

	if err != nil {
		log.Errorf("Unable connect to mongodb addr %s", host)
		return nil, err
	}
	// Only log the warning severity or above.
	//logrus.SetLevel(logrus.WarnLevel)
	log.Infof("connect to mongodb addr: %s", host)
	// Set the session mode as Eventual to ensure that the socket is created for each request.
	// Can switch to other mode only after the old APIs are cleaned up.
	session.SetMode(mgo.Eventual, true)

	go backgroundMongo(closing)

	err = ensureIndexes()
	if err != nil {
		log.Errorf("Fail to create indexes as %v", err)
		return nil, err
	}

	return mclosed, nil
}

// ensureIndexes ensures the indexes for each collection.
func ensureIndexes() error {
	projectCollection := session.DB(defaultDBName).C(projectCollectionName)
	projectIndex := mgo.Index{Key: []string{"name"}, Unique: true}
	err := projectCollection.EnsureIndex(projectIndex)
	if err != nil {
		log.Errorf("fail to create index for project as %v", err)
		return err
	}
	return nil
}

// NewStore copy a mongo client session
func NewStore() *DataStore {
	s := session.Copy()
	return &DataStore{
		s:                 s,
		projectCollection: session.DB(defaultDBName).C(projectCollectionName),
	}
}

// Close close mongo client session
func (d *DataStore) Close() {
	d.s.Close()
}

// Ping ping mongo server
func (d *DataStore) Ping() error {
	return d.s.Ping()
}

// Background goroutine for mongo. It can hold mongo connection & close session when progress exit.
func backgroundMongo(closing chan struct{}) {
	ticker := time.NewTicker(tickerDuration)
	for {
		select {
		case <-ticker.C:
			if err := session.Ping(); err != nil {
				log.Errorf("Ping Mongodb with error %s", err.Error())
				session.Refresh()
				session.SetSocketTimeout(socketTimeout)
				session.SetSyncTimeout(syncTimeout)
			}
		case <-closing:
			session.Close()
			log.Info("Mongodb session has been closed")
			close(mclosed)
			return
		}
	}
}
