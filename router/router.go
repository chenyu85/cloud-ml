package router

import (
	"cloud-ml/manager"
	"cloud-ml/model"
	"cloud-ml/store"
	"encoding/json"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	log "github.com/sirupsen/logrus"
)

const (
	// APIVersion is the version of API.
	APIVersion = "/api/v1"

	// projectPathParameterName represents the name of the path parameter for project.
	projectPathParameterName = "project"
)

// router represents the router to distribute the REST requests.
type router struct {
	// dataStore represents the manager for data store.
	dataStore *store.DataStore

	// projectManager represents the project manager.
	projectManager manager.ProjectManager
}

// InitRouters initializes the router for REST APIs.
func InitRouters(dataStore *store.DataStore) error {
	// New project manager
	projectManager, err := manager.NewProjectManager(dataStore)
	if err != nil {
		return err
	}

	router := &router{
		dataStore,
		projectManager,
	}

	ws := new(restful.WebService)
	ws.Filter(NCSACommonLogFormatLogger())
	router.registerProjectAPIs(ws)
	restful.Add(ws)
	return nil
}

// NCSACommonLogFormatLogger add filter to produce NCSA standard log.
func NCSACommonLogFormatLogger() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		t := time.Now()
		chain.ProcessFilter(req, resp)
		log.Infof("%d \"%s %s %s\" - %s %d %.2fs\n",
			resp.StatusCode(),
			req.Request.Method,
			req.Request.URL.RequestURI(),
			req.Request.Proto,
			req.Request.RemoteAddr,
			resp.ContentLength(),
			time.Since(t).Seconds(),
		)
	}
}

// registerProjectAPIs registers project related endpoints.
func (router *router) registerProjectAPIs(ws *restful.WebService) {
	log.Info("Register project APIs")
	ws.Path(APIVersion).Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)
	// POST /api/v1/projects
	ws.Route(ws.POST("/projects").To(router.createProject).
		Doc("Add a project").
		Reads(model.Project{}))

	// // GET /api/v1/projects
	// ws.Route(ws.GET("/projects").To(router.listProjects).
	// 	Doc("Get all projects"))

	// // PUT /api/v1/projects/{project}
	// ws.Route(ws.PUT("/projects/{project}").To(router.updateProject).
	// 	Doc("Update the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")).Reads(api.Project{}))

	// // GET /api/v1/projects/{project}
	// ws.Route(ws.GET("/projects/{project}").To(router.getProject).
	// 	Doc("Get the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")))

	// // DELETE /api/v1/projects/{project}
	// ws.Route(ws.DELETE("/projects/{project}").To(router.deleteProject).
	// 	Doc("Delete the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")))

	// // GET /api/v1/projects/{project}/repos
	// ws.Route(ws.GET("/projects/{project}/repos").To(router.listRepos).
	// 	Doc("List accessible repos of the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")))

	// // GET /api/v1/projects/{project}/branches
	// ws.Route(ws.GET("/projects/{project}/branches").To(router.listBranches).
	// 	Doc("List branches of the repo for the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")).
	// 	Param(ws.QueryParameter("repo", "the repo to list branches for").Required(true)))

	// // GET /api/v1/projects/{project}/tags
	// ws.Route(ws.GET("/projects/{project}/tags").To(router.listTags).
	// 	Doc("List tags of the repo for the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")).
	// 	Param(ws.QueryParameter("repo", "the repo to list branches for").Required(true)))

	// // GET /api/v1/projects/{project}/stats
	// ws.Route(ws.GET("/projects/{project}/stats").To(router.getProjectStatistics).
	// 	Doc("Get statistics of the project").
	// 	Param(ws.PathParameter("project", "name of the project").DataType("string")).
	// 	Param(ws.QueryParameter("startTime", "the start time of statistics").Required(false)).
	// 	Param(ws.QueryParameter("endTime", "the end time of statistics").Required(false)))
}

// EncodeResponse encodes response in json.
func EncodeResponse(rw http.ResponseWriter, statusCode int, data interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	return json.NewEncoder(rw).Encode(data)
}
