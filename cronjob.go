package main

import (
	"fmt"
	"log"
	"text/template"
)

func createCronjob(CronjobObj Cronjob) {
	log.Printf("Cronjobs for %s-%s-%s\n", application_name, CronjobObj.Name, build_id)
	fp := CreateFH("cronjobs.yaml")
	defer fp.Close()
	values := &Cronjobtmpl{
		Application_name:             application_name,
		Bamboo_deploy_release:        bamboo_deploy_release,
		Bamboo_AWS_HOSTNAME:          bamboo_AWS_HOSTNAME,
		Build_id:                     build_id,
		Build_nr:                     bamboo_buildNumber,
		CONSUL_APPLICATION:           CONSUL_APPLICATION,
		CONSUL_ENVIRONMENT:           CONSUL_ENVIRONMENT,
		CONSUL_FULL_URL:              CONSUL_FULL_URL,
		CONSUL_PASSWORD:              CONSUL_PASSWORD,
		CONSUL_URL:                   CONSUL_URL,
		CONSUL_USERNAME:              CONSUL_USERNAME,
        DEPLOYMENT_DATACENTER:        DEPLOYMENT_DATACENTER,
		Deploy:                       CronjobObj,
		Deploy_name:                  fmt.Sprintf("%s-%s-%s", application_name, CronjobObj.Name, build_id),
		Git_repo:                     git_repo,
		Namespace:                    deploy_namespace,
		NEW_RELIC_API_KEY_PASSWORD:   NEW_RELIC_API_KEY_PASSWORD,
		NEW_RELIC_API_URL:            NEW_RELIC_API_URL,
		NEW_RELIC_LICENSE_KEY:        NEW_RELIC_LICENSE_KEY,
        NEW_RELIC_ADMIN_KEY_PASSWORD: NEW_RELIC_ADMIN_KEY_PASSWORD,
		Ssh_key:                      ssh_key,
	}
	t := template.Must(template.ParseFiles("templates/cronjob.tmpl"))
	err := t.Execute(fp, values)
	if err != nil {
		log.Fatalf("ERROR: template execution: %s", err)
	}
	fp.Close()
}
