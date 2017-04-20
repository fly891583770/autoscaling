package routers

var MesosRouter = RouterArray{
	{Handler: "GetMarathonUrl", Method: "GET", Path: "/mesos/frameworks"},
}
