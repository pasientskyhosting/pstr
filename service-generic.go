package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func createGenericService(fp *os.File, AppObj App) {
	fmt.Printf("# Generic Service for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	values := &Servicetmpl{
		Application_name: application_name,
		Namespace:        deploy_namespace,
		Build_id:         build_id,
		Deploy:           AppObj,
	}

	t := template.Must(template.ParseFiles("templates/service-generic.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("template execution: %s", err)
		os.Exit(1)
	}
}
