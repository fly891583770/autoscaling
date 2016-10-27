package   main

type App_list struct{
	MarathonName                 string  `json:"marathon_name"`
}

type App_lnfo struct{
	MarathonName                string   `json:"marathon_name"`
    Appid                       string  `json:"app_id"`
}

type scaleUp struct{
	Mem					       float32   `json:"mem"`
	Cpu                        float32   `json:"cpu"`
	Thread                     float32   `json:"thread"`
	RequestQueue               float32   `json:"request_queue"`
	CollectPeriod              int       `json:"collect_period"`
	ContinuePeriod             int       `json:"continue_period"`
	CoolDownPeriod             int       `json:"cool_down_period"`
	ScaleAmount                int       `json:"scale_amount"`
	MaxScaleAmount             int       `json:"max_scale_amount"`
	Switch                     int       `json:"switch"`
} 


type scaleDown struct{
	Mem					     float32     `json:"mem"`
	Cpu                      float32     `json:"cpu"`
	Thread                   float32     `json:"thread"`
	RequestQueue             float32     `json:"request_queue"`
	CollectPeriod            int         `json:"collect_period"`
	ContinuePeriod           int         `json:"continue_period"`
	CoolDownPeriod           int         `json:"cool_down_period"`
	ScaleAmount              int         `json:"scale_amount"`
	MaxScaleAmount           int         `json:"max_scale_amount"`
	Switch                   int         `json:"switch"`
}

type CreateExpanionCap struct{
	MarathonName            string        `json:"marathon_name"`
	AppId                   string       `json:"app_id"`	
	ScaleUp                 scaleUp      `json:"scale_up"`
	ScaleDown               scaleDown    `json:"scale_down"` 
}


type UpdateExpaionCap  struct{
	MarathonName           string        `json:"marathon_name"`
	AppId                  string        `json:"app_id"`
	ScaleUp                scaleUp       `json:"scale_up"`
	ScaleDown              scaleDown     `json:"scale_down"`
}  

type Suspend struct{
	AppId                 string         `json:"app_id"`
	MarathonName           string        `json:"marathon_name"`	
}


type Start struct{
	AppId                string          `json:"app_id"`
	MarathonName           string        `json:"marathon_name"`
}


