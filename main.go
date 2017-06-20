package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	// Variables gotten from Environment
	bamboo_buildNumber         = os.Getenv("bamboo_buildNumber")
	bamboo_CONSUL_URL          = os.Getenv("bamboo_CONSUL_URL")
	bamboo_deploy_release      = os.Getenv("bamboo_deploy_release")
	bamboo_NEW_RELIC_API_URL   = os.Getenv("bamboo_NEW_RELIC_API_URL")
	cluster_ip                 = os.Getenv("cluster_ip")
	CONSUL_APPLICATION         = os.Getenv("bamboo_CONSUL_APPLICATION")
	CONSUL_ENVIRONMENT         = os.Getenv("bamboo_CONSUL_ENVIRONMENT")
	DEPLOYMENT_DATACENTER      = os.Getenv("bamboo_DEPLOYMENT_DATACENTER")
	CONSUL_PASSWORD            = os.Getenv("bamboo_CONSUL_PASSWORD")
	CONSUL_USERNAME            = os.Getenv("bamboo_CONSUL_USERNAME")
	NEW_RELIC_API_KEY_PASSWORD = os.Getenv("bamboo_NEW_RELIC_API_KEY_PASSWORD")
	NEW_RELIC_LICENSE_KEY      = os.Getenv("bamboo_NEW_RELIC_LICENSE_KEY_PASSWORD")
	ssh_key                    = os.Getenv("ssh_key")
	bamboo_AWS_HOSTNAME        = os.Getenv("bamboo_AWS_HOSTNAME")

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
	M_CLINIC          bool
	M_DEPLOY          bool
	M_GERNICSERVICE   bool
	M_INGRESS         bool
	M_SERVICE         bool
	NEW_RELIC_API_URL *url.URL
	O_FILENAME        string
	O_LIMIT           string
	O_OUTPUT          string
	clinic_name       string
	clinic_hostname   string
	D_HOSTNAMES       string
)

func checkErr(err error) {
	if err != nil {
		log.Printf("WARN: %#v", err)
	}
}

func init() {
	var err error
	CONSUL_URL, err = url.Parse(bamboo_CONSUL_URL)
	NEW_RELIC_API_URL, err = url.Parse(bamboo_NEW_RELIC_API_URL)
	if CONSUL_URL.Host != "" {
		CONSUL_FULL_URL, err = url.Parse(fmt.Sprintf("%s://%s:%s@%s/", CONSUL_URL.Scheme, CONSUL_USERNAME, CONSUL_PASSWORD, CONSUL_URL.Host))
	}
	git_repo, err = url.Parse(os.Getenv("git_repo"))

	_ = err
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
	flag.StringVar(&clinic_hostname, "clinic_hostname", "", "Clinic hostname")
	flag.StringVar(&clinic_name, "clinic_name", "", "Clinic name")
	flag.StringVar(&D_HOSTNAMES, "hostname", "", "Hostnames for ingress. comma separated")
	flag.StringVar(&deploy_namespace, "namespace", "", "namespace")
	flag.StringVar(&O_FILENAME, "file", "", "Filename to parse")
	flag.StringVar(&O_LIMIT, "limit", "", "Limit the run to certain app name")
	flag.StringVar(&O_OUTPUT, "output", "./", "Output folder")
	flag.Parse()

	if M_ALL {
		M_DEPLOY = true
		M_SERVICE = true
		M_GERNICSERVICE = true
		M_AUTOSCALER = true
		M_INGRESS = true
	}

	if clinic_hostname != "" {
		M_CLINIC = true
	}
}

func CreateFH(Filename string) (fp *os.File) {
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
	if M_CLINIC {
		createClinic(clinic_hostname)
	}

	if O_FILENAME != "" {
		hostnames = strings.Split(D_HOSTNAMES, ",")

		if build_id == "" || deploy_namespace == "" {
			log.Fatalf("Missing CMD line options build (\"%s\"), or namespace (\"%s\")", build_id, deploy_namespace)
		}
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

                // Set a default
                if AppObj.Scalability.TargetCPUUtilizationPercentage < 1 {
                    AppObj.Scalability.TargetCPUUtilizationPercentage = 70
                }

				if M_AUTOSCALER {
					if AppObj.Resources.Requests.Cpu != "" && AppObj.Scalability.MinReplicas > 0 && AppObj.Scalability.MaxReplicas > 1 && AppObj.Scalability.TargetCPUUtilizationPercentage > 0 {
						createAutoScaler(AppObj)
					}
				}
				if M_INGRESS {
					createIngress(AppObj)
				}

			}
			return true // keep iterating
		})
	} // end if O_FILENAME
} // end main
