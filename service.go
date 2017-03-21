package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func createService(Name string, AppObj App) {
	fmt.Printf("# Service for %s-%s-%s\n", application_name, Name, build_id)
	values := &Servicetmpl{
		Application_name: application_name,
		Namespace:        deploy_namespace,
		Name:             Name,
		Build_id:         build_id,
		Deploy:           AppObj,
	}
	t := template.Must(template.ParseFiles("templates/service.tmpl"))
	err := t.Execute(os.Stdout, values)
	if err != nil {
		log.Fatalf("template execution: %s", err)
		os.Exit(1)
	}
}
