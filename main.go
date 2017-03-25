package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"
)

var (
	// Variables gotten from Environment
	bamboo_buildNumber         = os.Getenv("bamboo_buildNumber")
	bamboo_deploy_release      = os.Getenv("bamboo_deploy_release")
	cluster_ip                 = os.Getenv("cluster_ip")
	CONSUL_APPLICATION         = os.Getenv("bamboo_CONSUL_APPLICATION")
	CONSUL_ENVIRONMENT         = os.Getenv("bamboo_CONSUL_ENVIRONMENT")
	CONSUL_PASSWORD            = os.Getenv("bamboo_CONSUL_PASSWORD")
	CONSUL_USERNAME            = os.Getenv("bamboo_CONSUL_USERNAME")
	NEW_RELIC_API_KEY_PASSWORD = os.Getenv("bamboo_NEW_RELIC_API_KEY_PASSWORD")
	NEW_RELIC_LICENSE_KEY      = os.Getenv("bamboo_NEW_RELIC_LICENSE_KEY_PASSWORD")
	ssh_key                    = os.Getenv("ssh_key")

	// Static configs
	application_name  string
	build_id          string
	CONSUL_FULL_URL   *url.URL
	CONSUL_URL        *url.URL
	deploy_build      string
	deploy_namespace  string
	git_repo          *url.URL
	hostnames         []string
	M_ALL             bool
	M_AUTOSCALER      bool
	M_DEPLOY          bool
	M_GERNICSERVICE   bool
	M_INGRESS         bool
	M_SERVICE         bool
	NEW_RELIC_API_URL *url.URL
	O_FILENAME        string
	O_LIMIT           string
	O_OUTPUT          string
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	CONSUL_URL, err = url.Parse(os.Getenv("bamboo_CONSUL_URL"))
	checkErr(err)
	NEW_RELIC_API_URL, err = url.Parse(os.Getenv("bamboo_NEW_RELIC_API_URL"))
	checkErr(err)
	CONSUL_FULL_URL, err = url.Parse(fmt.Sprintf("%s://%s:%s@%s/", CONSUL_URL.Scheme, CONSUL_USERNAME, CONSUL_PASSWORD, CONSUL_URL.Host))
	checkErr(err)
	git_repo, err = url.Parse(os.Getenv("git_repo"))
	checkErr(err)
	_ = git_repo
	_ = NEW_RELIC_API_URL
	_ = CONSUL_FULL_URL
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
	flag.StringVar(&O_OUTPUT, "output", "./", "Output folder")
	var D_HOSTNAMES = flag.String("hostname", "", "Hostnames for ingress. comma separated")
	flag.Parse()

	hostnames = strings.Split(*D_HOSTNAMES, ",")
	if build_id == "" || deploy_namespace == "" {
		log.Fatalf("Missing CMD line options build (\"%s\"), or namespace (\"%s\")", build_id, deploy_namespace)
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
	//var err error
	PFilename := fmt.Sprintf("%s/%s", path.Clean(O_OUTPUT), Filename)
	fp, err := os.OpenFile(PFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("ERROR: create file: ", err)
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
	file, err := ioutil.ReadFile(O_FILENAME)
	if err != nil {
		log.Fatalf("ERROR: File error: %v\n", err)
	}

	//Get the application name from Json object
	application_name = gjson.GetBytes(file, "application").Str

	Services := gjson.GetBytes(file, "services")

	if Services.Index == 0 {
		log.Fatal("ERROR: Json decode error, no services found\n")
	}

	Services.ForEach(func(key, value gjson.Result) bool {
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
