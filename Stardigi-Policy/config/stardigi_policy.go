package config

type PolicyConfig struct {
	Db          DB         `json:"db"`
	Prometheuss Prometheus `json:"prometheus"`
	Mesoss      Mesos      `json:"mesos"`
	Marathons   Marathon   `json:"marathon"`
	Logs        Log        `json:"log"`
}

type DB struct {
	Type       string `json:"type"`
	Datasource string `json:"datasource"`
}

type Prometheus struct {
	PrometheusAddr  string `json:"promethus_addr" `
	PrometheusPort  string `json:"promethus_port"`
	PrometheusGroup string `json:"promethus_group"`
	PrometheusJob   string `json:"promethus_job"`
}

type Mesos struct {
	MesosUrls string `json:"mesos_urls"`
}

type Marathon struct {
	MarathonName string `json:"marathon_name"`
}

type Log struct {
	Level string `json:"level"`
	File  string `json:"file"`
}
