package gitlabcli

type GitlabVarData struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description any    `json:"description"`
	Env         string `json:"environment_scope"`
	IsRaw       bool   `json:"raw"`
	IsHidden    bool   `json:"hidden"`
	IsProtected bool   `json:"protected"`
	IsMasked    bool   `json:"masked"`
}
