package httpc

type Cmth struct {
	Stauses string `json:"status"`
	Datas   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Results    []Result `json:"result"`
}

type Result struct {
	Metrics Metric        `json:"metric"`
	Values  []interface{} `json:"value"`
}

type Metric struct {
	MarathonAppId string `json:"marathon_app_id"`
	Name          string `json:"__name__"`
	AppId         string `json:"appId"`
	ContainerId   string `json:"containerId"`
	Group         string `json:"group"`
	Instance      string `json:"instance"`
	Job           string `json:"dcos"`
	MarathonName  string `json:"marathon"`
}
