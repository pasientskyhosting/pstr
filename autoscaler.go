package main

import (
	"fmt"
	"log"
	"text/template"
)

func createAutoScaler(AppObj App) {
	log.Printf("AutoScaler for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	fp := CreateFH("autoscaler.yaml")
	defer fp.Close()
	values := &Autoscalertmpl{
		Deploy:      AppObj,
		Deploy_name: fmt.Sprintf("%s-%s-%s", application_name, AppObj.Name, build_id),
		Namespace:   deploy_namespace,
	}
	t := template.Must(template.ParseFiles("templates/autoscaler.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("ERROR: template execution: %s", err)
	}
	fp.Close()
}
