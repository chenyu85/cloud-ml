package model

type Project struct {
	ID          string `bson:"_id,omitempty" json:"id,omitempty" description:"id of the project"`
	Name        string `bson:"name,omitempty" json:"name,omitempty" description:"name of the project, should be unique"`
	Description string `bson:"description,omitempty" json:"description,omitempty" description:"description of the project"`
}
