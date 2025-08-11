package gitlabcli

// type GitlabEnvDataNamespace struct {
//	Id int `json:"id"`
//}

// type GitlabEnvDataProject struct {
//	Id        int                `json:"id"`
//	Namespace GitlabEnvDataNamespace `json:"namespace"`
//}

type GitlabEnvData struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	State       any    `json:"state"`
	Url         any    `json:"external_url"`
	Description any    `json:"description"`
	//Project     GitlabEnvDataProject `json:"project"`
}
