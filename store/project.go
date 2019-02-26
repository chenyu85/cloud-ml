package store

import (
	"cloud-ml/model"

	"gopkg.in/mgo.v2/bson"
)

// CreateProject creates the project, returns the project created.
func (d *DataStore) CreateProject(project *model.Project) (*model.Project, error) {
	project.ID = bson.NewObjectId().Hex()
	if err := d.projectCollection.Insert(project); err != nil {
		return nil, err
	}
	return project, nil
}
