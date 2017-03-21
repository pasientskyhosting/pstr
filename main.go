package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Variables gotten from Environment
var bamboo_buildNumber = os.Getenv("bamboo_buildNumber")
var bamboo_deploy_release = os.Getenv("bamboo_deploy_release")
var build_id = os.Getenv("build_id")
var cluster_ip = os.Getenv("cluster_ip")
var CONSUL_APPLICATION = os.Getenv("CONSUL_APPLICATION")
var CONSUL_ENVIRONMENT = os.Getenv("CONSUL_ENVIRONMENT")
var CONSUL_PASSWORD = os.Getenv("CONSUL_PASSWORD")
var CONSUL_URL = os.Getenv("CONSUL_URL")
var CONSUL_USERNAME = os.Getenv("CONSUL_USERNAME")
var deploy_build string
var NEW_RELIC_LICENSE_KEY = os.Getenv("NEW_RELIC_LICENSE_KEY")
var ssh_key = os.Getenv("ssh_key")

// Static configs
var application_name string
var deploy_namespace string
var git_repo = os.Getenv("git_repo")
var hostnames []string
var M_ALL bool
var M_AUTOSCALER bool
var M_DEPLOY bool
var M_INGRESS bool
var M_SERVICE bool
var O_LIMIT string
var O_FILENAME string

func init() {
	flag.BoolVar(&M_ALL, "all", false, "Outputs deploymen, service, autoscaler and ingress")
	flag.BoolVar(&M_AUTOSCALER, "autoscaler", false, "Create autoscaler")
	flag.BoolVar(&M_DEPLOY, "deploy", false, "Create deployments")
	flag.BoolVar(&M_INGRESS, "ingress", false, "Create ingress rules")
	flag.BoolVar(&M_SERVICE, "service", false, "Create services")
	flag.StringVar(&deploy_build, "build", "", "build")
	flag.StringVar(&deploy_namespace, "namespace", "", "namespace for deployment")
	flag.StringVar(&O_LIMIT, "limit", "", "Limit the run to certain app name")
	flag.StringVar(&O_FILENAME, "file", "./serviceDefinition.json", "Filename to parse")
	var D_HOSTNAMES = flag.String("hostname", "", "Hostnames for ingress. comma separated")
	flag.Parse()
	if deploy_build == "" || deploy_namespace == "" {
		//if deploy_build == "" || deploy_namespace == "" || *D_HOSTNAMES == "" {
		println(deploy_build)
		println(deploy_namespace)
		println(hostnames)
		log.Fatal("Missing CMD line options build, or namespace")
		os.Exit(1)
	}

	hostnames = strings.Split(*D_HOSTNAMES, ",")
	if M_ALL {
		M_DEPLOY = true
		M_SERVICE = true
		M_AUTOSCALER = true
		M_INGRESS = true
	}
}

func Check_if_limit(AppObj App) bool {
	if len(O_LIMIT) > 0 {
		if AppObj.Name == O_LIMIT {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func main() {
	file, e := ioutil.ReadFile(O_FILENAME)
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
		} else if Check_if_limit(AppObj) {
			if M_DEPLOY {
				createDeploy(AppObj)
			}
			if M_SERVICE {
				createService(AppObj)
			}
			if M_AUTOSCALER {
				createAutoScaler(AppObj)
			}
			if M_INGRESS {
				createIngress(AppObj)
			}
		}
		return true // keep iterating
	})
}
