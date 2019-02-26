package manager

import (
	"cloud-ml/model"
	"cloud-ml/store"
	"fmt"

	"github.com/minio/minio-go"
	log "github.com/sirupsen/logrus"
)

type ProjectManager interface {
	CreateProject(project *model.Project) (*model.Project, error)
	// GetProject(projectName string) (*model.Project, error)
	// UpdateProject(projectName string, newProject *model.Project) (*model.Project, error)
	// DeleteProject(projectName string) error
}

// projectManager represents the manager for project.
type projectManager struct {
	dataStore *store.DataStore
}

// NewProjectManager creates a project manager.
func NewProjectManager(dataStore *store.DataStore) (ProjectManager, error) {
	if dataStore == nil {
		return nil, fmt.Errorf("Fail to new project manager as data store is nil.")
	}
	return &projectManager{dataStore}, nil
}

// CreateProject creates a project.
func (m *projectManager) CreateProject(project *model.Project) (*model.Project, error) {
	UploadFile()
	return nil, nil
}

func UploadFile() {
	endpoint := "localhost:9000"
	accessKeyID := "0KNYP5DNKSE18YUC6JLG"
	secretAccessKey := "W+luy3BAW9VaK21oWSkQGJEvDjSbNANLCCEDzOS0"
	useSSL := false

	// 初使化 minio client对象。
	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(minioClient)
	// 创建一个叫mymusic的存储桶。
	bucketName := "cloud"
	location := "us-east-1"

	err = minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n")
			log.Fatalln(err)
		}
	}
	log.Printf("Successfully created %s\n", bucketName)

	// 上传一个zip文件。
	objectName := "model.zip"
	filePath := "D:\\uploads\\model.zip"
	contentType := "application/zip"

	// 使用FPutObject上传一个zip文件。
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}
