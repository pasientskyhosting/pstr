package main

import (
	"log"
	"text/template"
)

func createService(AppObj App) {
	log.Printf("# Service for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	fp := CreateFH("service.yaml")
	defer fp.Close()
	values := &Servicetmpl{
		Application_name: application_name,
		Namespace:        deploy_namespace,
		Build_id:         build_id,
		Deploy:           AppObj,
	}

	t := template.Must(template.ParseFiles("templates/service.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("ERROR: template execution: %s", err)
	}
	fp.Close()
}
