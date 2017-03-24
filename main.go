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

var (
	// Variables gotten from Environment
	bamboo_buildNumber    = os.Getenv("bamboo_buildNumber")
	bamboo_deploy_release = os.Getenv("bamboo_deploy_release")

	//var build_id = os.Getenv("build_id")
	cluster_ip                 = os.Getenv("cluster_ip")
	CONSUL_APPLICATION         = os.Getenv("bamboo_CONSUL_APPLICATION")
	CONSUL_ENVIRONMENT         = os.Getenv("bamboo_CONSUL_ENVIRONMENT")
	CONSUL_PASSWORD            = os.Getenv("bamboo_CONSUL_PASSWORD")
	CONSUL_URL                 = os.Getenv("bamboo_CONSUL_URL")
	CONSUL_USERNAME            = os.Getenv("bamboo_CONSUL_USERNAME")
    CONSUL_FULL_URL            = "https://" + CONSUL_USERNAME + ":" + CONSUL_PASSWORD + "@" + CONSUL_URL
	NEW_RELIC_LICENSE_KEY      = os.Getenv("bamboo_NEW_RELIC_LICENSE_KEY_PASSWORD")
	NEW_RELIC_API_URL          = os.Getenv("bamboo_NEW_RELIC_API_URL")
	NEW_RELIC_API_KEY_PASSWORD = os.Getenv("bamboo_NEW_RELIC_API_KEY_PASSWORD")
	ssh_key                    = os.Getenv("ssh_key")
	git_repo                   = os.Getenv("git_repo")

	// Static configs
	deploy_build     string
	application_name string
	deploy_namespace string
	build_id         string
	hostnames        []string
	M_ALL            bool
	M_AUTOSCALER     bool
	M_DEPLOY         bool
	M_INGRESS        bool
	M_SERVICE        bool
	M_GERNICSERVICE  bool
	O_LIMIT          string
	O_FILENAME       string
	O_OUTPUT         string
)

func init() {
	flag.BoolVar(&M_ALL, "all", false, "Outputs deploymen, service, autoscaler and ingress")
	flag.BoolVar(&M_AUTOSCALER, "autoscaler", false, "Create autoscaler")
	flag.BoolVar(&M_DEPLOY, "deploy", false, "Create deployments")
	flag.BoolVar(&M_INGRESS, "ingress", false, "Create ingress rules")
	flag.BoolVar(&M_SERVICE, "service", false, "Create services")
	flag.BoolVar(&M_GERNICSERVICE, "genericservice", false, "Create generic services")
	flag.StringVar(&build_id, "build_id", "", "build_id from bamboo")
	flag.StringVar(&deploy_namespace, "namespace", "", "namespace for deployment")
	flag.StringVar(&O_LIMIT, "limit", "", "Limit the run to certain app name")
	flag.StringVar(&O_FILENAME, "file", "serviceDefinition.json", "Filename to parse")
	flag.StringVar(&O_OUTPUT, "output", "", "Output folder")
	var D_HOSTNAMES = flag.String("hostname", "", "Hostnames for ingress. comma separated")
	flag.Parse()
	hostnames = strings.Split(*D_HOSTNAMES, ",")

	if build_id == "" || deploy_namespace == "" {
		//if deploy_build == "" || deploy_namespace == "" || *D_HOSTNAMES == "" {
		println(deploy_build)
		println(deploy_namespace)
		println(hostnames)
		log.Fatal("Missing CMD line options build, or namespace")
		os.Exit(1)
	}

	if M_ALL {
		M_DEPLOY = true
		M_SERVICE = true
		M_GERNICSERVICE = true
		M_AUTOSCALER = true
		M_INGRESS = true
	}
}

func CreateFH(Filename string) (fp *os.File) {
	var err error
	if O_OUTPUT != "" {
		fp, err = os.OpenFile(Filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Println("create file: ", err)
			os.Exit(1)
		}
	} else {
		fp = os.Stdout
	}
	return fp
}

func Check_if_limit(AppObj App) bool {
	if O_LIMIT != "" {
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
			if M_GERNICSERVICE {
				createGenericService(AppObj)
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
