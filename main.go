package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/cloudfoundry-community/go-cfenv"
	"github.com/gorilla/mux"
)

type reqHandler struct {
	cfClient *cfclient.Client
}

type synchReqBody struct {
	OrgName string `json:"orgname"`
}

type synchReqResponseBody struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (h reqHandler) synchOrg(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var requestBody synchReqBody

	responseBody := synchReqResponseBody{
		Success: true,
		Message: "synch done",
	}

	w.Header().Set("content-type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	org, err := h.cfClient.GetOrgByName(requestBody.OrgName)
	if err != nil {
		responseBody.Success = false
		responseBody.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(responseBody)
		return
	}

	log.Printf("org found with name: %v, GUID: %v\n", org.Name, org.Guid)

	secGroup, err := h.cfClient.GetSecGroupByName(org.Name)
	if err != nil {
		log.Println(err.Error())
		responseBody.Success = false
		responseBody.Message = err.Error()
		json.NewEncoder(w).Encode(responseBody)
		return
	}

	log.Printf("ASG found with name: %v, GUID: %v\n", secGroup.Name, secGroup.Guid)

	spaces, _ := h.cfClient.OrgSpaces(org.Guid)
	for _, space := range spaces {
		log.Printf("Binding space %v to ASG %v", space.Name, secGroup.Name)
		err := h.cfClient.BindSecGroup(secGroup.Guid, space.Guid)
		if err != nil {
			log.Printf("Failed binding asg to space. Error: %v", err.Error())
		}
	}

	json.NewEncoder(w).Encode(responseBody)
}

func main() {

	appEnv, _ := cfenv.Current()
	user := os.Getenv("CFUSER")
	password := os.Getenv("CFPASSWORD")

	if user == "" || password == "" {
		log.Fatal("Please provide CFUSER and CFPASSWORD environment variables")
	}

	c := &cfclient.Config{
		ApiAddress: appEnv.CFAPI,
		Username:   user,
		Password:   password,
	}

	log.Println("Connecting to CF API")
	client, _ := cfclient.NewClient(c)

	synchHandler := reqHandler{
		cfClient: client,
	}

	log.Println("Connected")

	log.Println("Starting http server")
	r := mux.NewRouter()
	r.Path("/api/v1/synch").Methods(http.MethodPost).HandlerFunc(synchHandler.synchOrg)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", appEnv.Port), r))
}
