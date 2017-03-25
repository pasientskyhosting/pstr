package main

import (
	"log"
	"text/template"
)

func createGenericService(AppObj App) {
	log.Printf("# Generic Service for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	fp := CreateFH("service-generic.yaml")
	defer fp.Close()
	values := &Servicetmpl{
		Application_name: application_name,
		Build_id:         build_id,
		Deploy:           AppObj,
		Namespace:        deploy_namespace,
	}

	t := template.Must(template.ParseFiles("templates/service-generic.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	fp.Close()
}
