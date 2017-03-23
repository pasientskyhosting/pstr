package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func createDeploy(fp *os.File, AppObj App) {
	fmt.Printf("# Deployment for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	values := &Deploytmpl{
		Application_name:      application_name,
		Bamboo_deploy_release: bamboo_deploy_release,
		Build_id:              build_id,
		Build_nr:              bamboo_buildNumber,
		CONSUL_APPLICATION:    CONSUL_APPLICATION,
		CONSUL_PASSWORD:       CONSUL_PASSWORD,
		CONSUL_URL:            CONSUL_URL,
		CONSUL_USERNAME:       CONSUL_USERNAME,
		CONSUL_ENVIRONMENT:    CONSUL_ENVIRONMENT,
		Deploy:                AppObj,
		Deploy_name:           application_name + "-" + AppObj.Name + "-" + build_id,
		Git_repo:              git_repo,
		Namespace:             deploy_namespace,
		NEW_RELIC_LICENSE_KEY: NEW_RELIC_LICENSE_KEY,
        NEW_RELIC_API_URL:     NEW_RELIC_API_URL,
        NEW_RELIC_API_KEY_PASSWORD: NEW_RELIC_API_KEY_PASSWORD,
		Ssh_key:               ssh_key,
	}

	t := template.Must(template.ParseFiles("templates/deploy.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("template execution: %s", err)
		os.Exit(1)
	}
}
