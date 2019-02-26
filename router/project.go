package router

import (
	"cloud-ml/model"
	"log"
	"net/http"

	"github.com/emicklei/go-restful"
)

// createProject handles the request to create a project.
func (router *router) createProject(request *restful.Request, response *restful.Response) {
	project := &model.Project{}
	createdProject, err := router.projectManager.CreateProject(project)
	if err != nil {
		log.Fatalf("create Project failed,err msg%s", err)
	}
	response.WriteHeaderAndEntity(http.StatusCreated, createdProject)
}
