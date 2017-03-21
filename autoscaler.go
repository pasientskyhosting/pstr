package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func createAutoScaler(Name string, AppObj App) {
	fmt.Printf("# AutoScaler for %s-%s-%s\n", application_name, Name, build_id)
	values := &Autoscalertmpl{
		Deploy_name: application_name + "-" + Name + "-" + build_id,
		Namespace:   deploy_namespace,
		Name:        Name,
		Deploy:      AppObj,
	}

	t := template.Must(template.ParseFiles("templates/autoscaler.tmpl"))
	err := t.Execute(os.Stdout, values)
	if err != nil {
		log.Fatalf("template execution: %s", err)
		os.Exit(1)
	}
}
