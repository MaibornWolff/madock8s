package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GreetingRequest struct {
	Name string `json:"name"`
}

type GreetingResponse struct {
	Result struct {
		Greeting string `json:"greeting"`
	} `json:"result"`
}

func GreetingHandler(w http.ResponseWriter, req *http.Request) {

	greetingRequest := GreetingRequest{}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	json.Unmarshal(body, &greetingRequest)

	resp := GreetingResponse{}

	resp.Result.Greeting = "Welcome to maDocK8s, developer " + greetingRequest.Name

	jsonResp, _ := json.Marshal(resp)
	fmt.Fprint(w, string(jsonResp))
}
