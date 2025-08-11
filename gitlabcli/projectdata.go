package gitlabcli

type GitlabProjectData struct {
	Id                int    `json:"id"`
	Name              string `json:"name"`
	Description       any    `json:"description"`
	Path              string `json:"path"`
	NameWithNamespace string `json:"name_with_namespace"`
	PathWithNamespace string `json:"path_with_namespace"`
	SshUrlToRepo      string `json:"ssh_url_to_repo"`
	HttpUrlToRepo     string `json:"http_url_to_repo"`
	WebUrl            string `json:"web_url"`
	Visibility        string `json:"visibility"`
}
