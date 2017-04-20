package httpc

import (
	"Stardigi-Policy/routers"
	"Stardigi-Policy/utils"
	"fmt"
	"io"
	// "io/ioutil"
	"net/http"
	"strings"
	"time"
)

type StarRequestGen struct {
	Host    string
	Routers routers.RouterArray
}

func NewStarRequestGen(host string, router routers.RouterArray) *StarRequestGen {

	// 添加路由格式检查

	return &StarRequestGen{
		Host:    host,
		Routers: router,
	}
}

func createHttpClient(timeout time.Duration) *http.Client {

	client := &http.Client{
		Timeout: timeout,
	}

	return client
}

/**
 * the function declare SSH transmit
 */
func CreateHttpsClient() {

}

/**
 * @param handler
 * @param param : param of url ,such as /:name -> /ylc
 * @param body
 * @param fParam : param of table ,such as /getname?xxx
 */
func (star *StarRequestGen) CreateStarRequest(handler string, param Mapstring, body io.Reader, header Mapstring, fParam string) (*http.Request, error) {
	// fParam is a form param of url
	router := star.LookUrl(handler)
	path := router.Path
	//设置路由参数
	if len(param) > 0 {
		path = utils.ParaseUrlParam(router.Path, param)
	}

	endpoint := utils.StringJoin(star.Host, path)

	// 设置表单参数
	if !strings.EqualFold(fParam, "") {
		endpoint = utils.StringJoin(endpoint, "?query=", fParam)
	}

	fmt.Println("===发送HTTP请求的地址是==", endpoint)

	request, err := http.NewRequest(router.Method, endpoint, body)

	if err != nil {
		return request, err
	}

	// 设置请求头
	if header != nil {
		for key, value := range header {
			request.Header.Set(key, value)
		}
	}

	return request, nil
}

func (star *StarRequestGen) LookUrl(handler string) *routers.Router {

	router := new(routers.Router)

	for _, r := range star.Routers {

		if strings.EqualFold(r.Handler, handler) {
			*router = r
		}
	}

	return router
}

func (star *StarRequestGen) DoHttpRequest(handler string, param Mapstring, body io.Reader, header Mapstring, fParam string) (chan *http.Response, chan error) {

	responseChan := make(chan *http.Response, 1)
	errChan := make(chan error, 1)
	timeout := 100 * time.Second
	client := createHttpClient(timeout)

	request, err := star.CreateStarRequest(handler, param, body, header, fParam)
	if err != nil {
		errChan <- err
	}
	response, err := client.Do(request)
	if err != nil {
		errChan <- err
	}
	if response != nil {
		responseChan <- response
	} else {
		fmt.Println("=====response为空====")
	}
	return responseChan, errChan
}
