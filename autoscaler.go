package main

import (
	"fmt"
	"log"
	//"os"
	"text/template"
)

func createAutoScaler(AppObj App) {
	fp := CreateFH(O_OUTPUT + "/autoscaler.yaml")
	fmt.Printf("# AutoScaler for %s-%s-%s\n", application_name, AppObj.Name, build_id)
	values := &Autoscalertmpl{
		Deploy_name: application_name + "-" + AppObj.Name + "-" + build_id,
		Namespace:   deploy_namespace,
		Deploy:      AppObj,
	}
	t := template.Must(template.ParseFiles("templates/autoscaler.tmpl"))
	err := t.Execute(fp, values)

	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
	fp.Close()
}
