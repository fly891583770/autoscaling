package routers

var MarathonRouter = RouterArray{
	{Handler: "UpdateApp", Method: "PUT", Path: "/v2/apps/:app_id"},
	{Handler: "GetApps", Method: "GET", Path: "/v2/apps/:app_id"},
}
