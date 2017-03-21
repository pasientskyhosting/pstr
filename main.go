package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
)

var deploy_name string
var deploy_build string
var deploy_namespace string
var application_name string
var bamboo_deploy_release = os.Getenv("bamboo_deploy_release")
var build_id = os.Getenv("build_id")
var bamboo_buildNumber = os.Getenv("bamboo_buildNumber")
var CONSUL_APPLICATION = os.Getenv("CONSUL_APPLICATION")
var CONSUL_ENVIRONMENT = os.Getenv("CONSUL_ENVIRONMENT")
var CONSUL_PASSWORD = os.Getenv("CONSUL_PASSWORD")
var CONSUL_URL = os.Getenv("CONSUL_URL")
var CONSUL_USERNAME = os.Getenv("CONSUL_USERNAME")
var cluster_ip = os.Getenv("cluster_ip")
var git_repo = os.Getenv("git_repo")
var ssh_key = os.Getenv("ssh_key")
var NEW_RELIC_LICENSE_KEY = os.Getenv("NEW_RELIC_LICENSE_KEY")

func init() {
	var D_NAME = flag.String("name", "", "name")
	var D_BUILD = flag.String("build", "", "build")
	var D_NAMESPACE = flag.String("namespace", "", "namespace")

	flag.Parse()
	if *D_NAME == "" || *D_BUILD == "" || *D_NAMESPACE == "" {
		os.Exit(1)
	}
	deploy_name = *D_NAME
	deploy_build = *D_BUILD
	deploy_namespace = *D_NAMESPACE

}

func main() {
	file, e := ioutil.ReadFile("./deploy.json")
	if e != nil {
		fmt.Fprintf(os.Stderr, "File error: %v\n", e)
		os.Exit(1)
	}
	//Get the application name from Json object
	application_name = gjson.GetBytes(file, "application").Str
	
	value := gjson.GetBytes(file, "services")

	if value.Index == 0 {
		fmt.Fprint(os.Stderr, "Json decode error, no services found\n")
		os.Exit(1)
	}

	value.ForEach(func(key, value gjson.Result) bool {

		var AppObj App
		err := json.Unmarshal([]byte(value.String()), &AppObj)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Json Error: %s\n", err)
			os.Exit(1)
		} else {
			// fmt.Printf("# Start %s-%s\n", application_name, key.String())
			createDeploy(key.String(), AppObj)
			createAutoScaler(key.String(), AppObj)
			createService(key.String(), AppObj)
			createIngress(key.String(), AppObj)
			// fmt.Printf("# End %s-%s\n\n", application_name, key.String())
		}
		return true // keep iterating
	})
}
