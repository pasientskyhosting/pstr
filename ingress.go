package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func createIngress(AppObj App) {
	if AppObj.Ports.External.HTTP > 0 {

		fmt.Printf("# Ingress for %s-%s-%s\n", application_name, AppObj.Name, build_id)
		values := &Ingresstmpl{
			Application_name: application_name,
			Namespace:        deploy_namespace,
			Deploy:           AppObj,
			Hostnames:        hostnames,
		}

		t := template.Must(template.ParseFiles("templates/ingress.tmpl"))
		err := t.Execute(os.Stdout, values)
		if err != nil {
			log.Fatalf("template execution: %s", err)
			os.Exit(1)
		}

	} else {
		fmt.Fprintf(os.Stderr, "# %s-%s has no external port, Skipping Ingress\n", application_name, AppObj.Name)
	}
}
