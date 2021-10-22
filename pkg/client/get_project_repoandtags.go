package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type RepoInfo struct {
	Name string `json:"name"`
	ID uint `json:"id"`
	//ProjectId uint `json:"project_id"`
}

type TagInfo struct {
	Name string `json:"name"`
}

type ProjectInfo struct {
	Name string `json:"name"`
	ProjectId int `json:"project_id"`
}

// 获取指定repo 的所有tag
func getRepositoryTag(oldRegistry string,projectName string, repoName string, oldRegistryUser string,oldRegistryPwd string) (tags []string) {
	timeout := 10 * time.Second //超时时间50ms
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	//reqUrl := "https://" + Registry + "/api/repositories/" + repoName + "/tags"
	reqUrl := fmt.Sprintf("https://%s/api/repositories/%s/%s/tags", oldRegistry, projectName, repoName)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(oldRegistryUser, oldRegistryPwd)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var tagInfo []TagInfo
	json.Unmarshal(body, &tagInfo)
	for _, v := range tagInfo {
		tags = append(tags, v.Name)
	}
	return tags
}

// 从 项目名称获取对应的项目id，如果查询不到，则返回 -1
func getProjectIdFromName(oldRegistry string, projectName string,oldRegistryUser string,oldRegistryPwd string) (projectId int) {
	timeout := 10 * time.Second //超时时间50ms
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	//reqUrl := "https://" + Registry + "/api/repositories/" + repoName + "/tags"
	reqUrl := fmt.Sprintf("https://%s/api/projects", oldRegistry)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(oldRegistryUser, oldRegistryPwd)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	var projectInfos []ProjectInfo
	json.Unmarshal(body, &projectInfos)
	//fmt.Println(projectInfos)
	for _, v := range projectInfos {
		// 如果projectName 符合传入的预期，则直接返回
		if v.Name == projectName {
			return v.ProjectId
		}
	}
	return -1
}


// 获取指定 project id 的所有 repository
func getRepositories(oldRegistry string, projectId uint, oldRegistryUser string, oldRegistryPwd string) (repoInfo []RepoInfo) {
	timeout := 10 * time.Second //超时时间50ms
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	// 获取 project 的所有 repository 的 api路径
	reqUrl := fmt.Sprintf("https://%s/api/repositories?project_id=%d", oldRegistry, projectId)
	//reqUrl := "https://" + Registry + "/api/repositories?project_id=7"
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(oldRegistryUser, oldRegistryPwd)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &repoInfo)

	// 把 public/busybox 变成 busybox
	for k, v := range repoInfo {
		repoInfo[k].Name = wipeRepoName(v.Name)
	}
	return repoInfo
}

func wipeRepoName(this string) string {
	i := strings.Index(this, "/")
	return this[i+1:]
}