package httpClient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/selefra/selefra/global"
	"io"
	"net/http"
	"os"
)

type OutputReq struct {
	Name     string                 `json:"name"`
	Query    string                 `json:"query"`
	Labels   map[string]interface{} `json:"labels"`
	Metadata []Metadata             `json:"metadata"`
}

type Metadata struct {
	Id                string `json:"id"`
	Severity          string `json:"severity"`
	Provider          string `json:"provider"`
	ResourceType      string `json:"resource_type"`
	ResourceAccountId string `json:"resource_account_id"`
	ResourceId        string `json:"resource_id"`
	ResourceRegion    string `json:"resource_region"`
	Remediation       string `yaml:"remediation" json:"remediation"`
	Title             string `json:"title"`
	Description       string `json:"description"`
	Output            string `json:"output"`
}

type OutputRes struct {
}

type UploadWorkplaceRes struct {
}

type Res[T any] struct {
	Code int    `json:"code"`
	Data T      `json:"data"`
	Msg  string `json:"msg"`
}

type CreateProjectData struct {
	Name    string `json:"name"`
	OrgName string `json:"org_name"`
	TaskId  uint   `json:"task_id"`
}

type loginData struct {
	UserName  string `json:"user_name"`
	TokenName string `json:"token_name"`
	OrgName   string `json:"org_name"`
}

type TaskData struct {
	TaskUUID string `json:"task_uuid"`
}
type logoutData struct {
}

type SetupStagRes struct{}

type dsnData struct {
	Dsn string `json:"dsn"`
}

type WorkPlaceReq struct {
	Data        []Data `json:"data"`
	ProjectName string `json:"project_name"`
	Token       string `json:"token"` // 示例Token:e0b82826a400250f68da273c22448a53
}

type Data struct {
	Path        string `json:"path"`
	YAMLContent string `json:"yaml_content"`
}

func CliHttpClient[T any](method, url string, info interface{}) (*Res[T], error) {
	var client http.Client
	bytesData, err := json.Marshal(info)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, "https://"+global.SERVER+url, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("404 not found")
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res Res[T]
	err = json.Unmarshal(respBytes, &res)
	if err != nil {
		return nil, err
	}
	return &res, err
}

func Login(token string) (*Res[loginData], error) {
	var info = make(map[string]string)
	info["token"] = token
	res, err := CliHttpClient[loginData]("POST", "/cli/login", info)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf(res.Msg)
	}
	return res, nil
}

func CreateTask(token, project_name string, task_id uint) (*Res[TaskData], error) {
	var info = make(map[string]interface{})
	info["token"] = token
	info["project_name"] = project_name
	info["task_id"] = task_id
	info["task_source"] = os.Getenv("SELEFRA_TASK_SOURCE")
	res, err := CliHttpClient[TaskData]("POST", "/cli/create_task", info)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return nil, fmt.Errorf(res.Msg)
	}
	return res, nil
}

func Logout(token string) error {
	var info = make(map[string]string)
	info["token"] = token
	res, err := CliHttpClient[logoutData]("POST", "/cli/logout", info)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}
	return nil
}

func CreateProject(token, name string) (orgName string, taskId uint, err error) {
	var info = make(map[string]string)
	info["token"] = token
	info["name"] = name
	res, err := CliHttpClient[CreateProjectData]("POST", "/cli/create_project", info)
	if err != nil {
		return "", 0, err
	}
	if res.Code != 0 {
		return "", 0, fmt.Errorf(res.Msg)
	}
	return res.Data.OrgName, res.Data.TaskId, nil
}
func GetDsn(token string) (string, error) {
	var info = make(map[string]string)
	info["token"] = token
	res, err := CliHttpClient[dsnData]("POST", "/cli/fetch_dsn", info)
	if err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf(res.Msg)
	}
	return res.Data.Dsn, nil
}

func OutPut(token, project string, req []OutputReq) error {
	var info = make(map[string]interface{})
	info["data"] = req
	info["token"] = token
	info["project_name"] = project
	res, err := CliHttpClient[OutputRes]("POST", "/cli/upload_issue", info)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}
	return nil
}

func UploadWorkplace(token, project string, fileMap map[string]string) error {

	var workplace WorkPlaceReq

	workplace.Token = token
	workplace.ProjectName = project
	workplace.Data = make([]Data, 0)
	for k, v := range fileMap {
		workplace.Data = append(workplace.Data, Data{
			Path:        k,
			YAMLContent: v,
		})
	}
	res, err := CliHttpClient[UploadWorkplaceRes]("POST", "/cli/upload_workplace", workplace)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return errors.New(res.Msg)
	}
	return nil
}

const Creating = "creating"
const Testing = "testing"
const Failed = "failed"

func SetupStag(token, project, stag string) error {
	if token == "" {
		return errors.New("token is empty")
	}
	var info = make(map[string]string)
	info["token"] = token
	info["project_name"] = project
	info["stag"] = stag
	res, err := CliHttpClient[SetupStagRes]("POST", "/cli/update_setup_stag", info)
	if err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf(res.Msg)
	}
	return nil
}
