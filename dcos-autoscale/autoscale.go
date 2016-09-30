package main

import (
	"fmt"
	"github.com/dcos/dcos-autoscale/common"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
)

func auto_query_list(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/applist" 
	fmt.Println("api url %s", new_query_list_Url)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.NewHttpError("auto_query_list failed", http.StatusInternalServerError)
	}

	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read auto_query_list resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}

func auto_query_details(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/appinfo" 
	fmt.Println("api url %s", new_query_list_Url)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return common.NewHttpError("auto_query_details failed", http.StatusInternalServerError)
	}
	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read auto_query_details resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}

func elastic_expansion_configuration(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/ruleset"
	fmt.Println("api url %s", new_query_list_Url)
    fmt.Println(r.Body)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return common.NewHttpError("elastic_expansion_configuration failed", http.StatusInternalServerError)
	}

	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read elastic_expansion_configuration resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}

func updating_elastic_expansion(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/ruleupdate"
	fmt.Println("api url %s", new_query_list_Url)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return common.NewHttpError("updating_elastic_expansion failed", http.StatusInternalServerError)
	}
	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read updating_elastic_expansion resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}

func suspend(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/pause" 
	fmt.Println("api url %s", new_query_list_Url)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.NewHttpError("suspend failed", http.StatusInternalServerError)
	}
	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read suspend resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}

func turn_on(ctx context.Context, w http.ResponseWriter, r *http.Request) *common.HttpError {

	client := &http.Client{}
	new_query_list_Url := ctx.Value("listaddr").(string) + "/recover" 
	fmt.Println("api url %s", new_query_list_Url)
	req, _ := http.NewRequest("POST", new_query_list_Url, r.Body)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return common.NewHttpError("turn_on failed", http.StatusInternalServerError)
	}
	list, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return common.NewHttpError("Read turn_on resp body failed", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(list)
	return nil

}
