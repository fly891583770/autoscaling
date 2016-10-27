package  main

type Applist struct{
	APP []AppList      		    `json:"apps"`
}

type AppList struct{
	AppId         		string  `json:"app_id"`
	Status         		string  `json:"status"`
	CountStatus    		string  `json:"count_status"`
	EeventDescription   string  `json:"event_description"`
}

type scaleup struct{
	Mem					float64 `json:"mem"`
	Cpu                 float64 `json:"cpu"`
	Thread              float64 `json:"thread"`
	RequestQueue        float64  `json:"request_queue"`
	CollectPeriod       int  	`json:"collect_period"`
	ContinuePeriod      int  	`json:"continue_period"`
	CoolDownPeriod      int 	`json:"cool_down_period"`
	ScaleAmount         int  	`json:"scale_amount"`
	MaxScaleAmount      int  	`json:"max_scale_amount"`
	Switch              int  	`json:"switch"`
} 

type scaledown struct{

	Mem					float64  `json:"mem"`
	Cpu                 float64  `json:"cpu"`
	Thread              float64  `json:"thread"`
	RequestQueue        float64  `json:"request_queue"`
	CollectPeriod       int      `json:"collect_period"`
	ContinuePeriod      int      `json:"continue_period"`
	CoolDownPeriod      int      `json:"cool_down_period"`
	ScaleAmount         int      `json:"scale_amount"`
	MaxScaleAmount      int      `json:"max_scale_amount"`
	Switch              int      `json:"switch"`
}

type QueryDetail struct{
	AppId         		string  `json:"app_id"`
	MarathonName        string  `json:"marathon_name"`
	ScaleUp             scaleup  `json:"scale_up"`
	ScaleDown           scaledown `json:"scale_down"`
}

type SuspendInfo struct{
	Status 				string  `json:"status"`
	Msgup               string  `json:"msgup"`
	Msgdown             string  `json:"msgdown"`
} 


type StartInfo struct{
	Status  			string  `json:"status"`
	Msgup     			string  `json:"msg"`
	Msgdown             string  `json:"msgdown"`
}

type SetExpanionCap struct{
	Status  			string  `json:"status"`
	Msgup     			string  `json:"msg"`
	Msgdown             string  `json:"msgdown"`
}

type UpExpanionCap struct{
	Status  			string  `json:"status"`
	Msgup     			string  `json:"msg"`
	Msgdown             string  `json:"msgdown"`
}

type ErrorStatus struct{
	Status  			string  `json:"status"`
}