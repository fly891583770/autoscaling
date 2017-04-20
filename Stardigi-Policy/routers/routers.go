package routers

type Router struct {
	Handler string
	Method  string
	Path    string
}

type RouterArray []Router
