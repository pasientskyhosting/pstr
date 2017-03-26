package main

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"
	"strings"
	"text/template"
)

func createClinic(clinic_hostnames string) {
	if deploy_namespace == "" {
		log.Fatal("Missing namespace to create clinic ingress")
	}
	c_hostnames := strings.Split(clinic_hostname, ",")
	fp := CreateFH("clinic-ingress.yaml")
	defer fp.Close()

	funcMap := template.FuncMap{
		"getDomain": func(FQDN string) string {
			res, err := publicsuffix.Domain(FQDN)
			if err != nil {
				log.Fatalf("ERROR: getDomain: %s", err)
			}
			return res
		},
	}

	for _, value := range c_hostnames {
		log.Printf("Ingress for %s\n", value)

		values := &Clinictmpl{
			Clinic_hostname: value,
			Namespace:       deploy_namespace,
		}
		t, _ := template.New("clinic-ingress.tmpl").Funcs(funcMap).ParseFiles("templates/clinic-ingress.tmpl")
		err := t.Execute(fp, values)

		if err != nil {
			log.Fatalf("ERROR: template execution: %s", err)
		}
	}

	fp.Close()
}
