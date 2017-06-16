package main

import (
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"log"
	"text/template"
    "strings"
)

func createIngress(AppObj App) {
	if AppObj.Ports.External.HTTP > 0 {
		log.Printf("Ingress for %s-%s-%s\n", application_name, AppObj.Name, build_id)
		fp := CreateFH("ingress.yaml")
		defer fp.Close()

		funcMap := template.FuncMap{
			"getDomain": func(FQDN string) string {
				res, err := publicsuffix.Domain(FQDN)
				if err != nil {
					log.Fatalf("ERROR: getDomain: %s", err)
				}
				return res
			},
            "mapDomainToCert": func (hostname string) string {
                // List of domains
                list := map[string]string {
                    // Subdomains - order matters
                    "svc.pasientsky.no": "star.svc.pasientsky.no",

                    // Domains
                    "patientsky.no": "star.patientsky.no",
                    "pasientsky.no": "star.pasientsky.no",
                    "patientsky.com": "star.patientsky.com",
                    "gel.camp": "star.gel.camp",
                    "publicdns.zone": "star.publicdns.zone",
                    "privatedns.zone": "star.privatedns.zone",
                }

                // Find cert from hostname
                for domain, cert := range list {
                    if strings.Contains(hostname, domain) {
                        return cert
                    }
                }

                // Return default domain
                return "star.publicdns.zone"
            },
		}

		values := &Ingresstmpl{
			Application_name: application_name,
			Build_id:         build_id,
			Deploy:           AppObj,
			Hostnames:        hostnames,
			Namespace:        deploy_namespace,
		}
		t, _ := template.New("ingress.tmpl").Funcs(funcMap).ParseFiles("templates/ingress.tmpl")
		err := t.Execute(fp, values)

		if err != nil {
			log.Fatalf("ERROR: template execution: %s", err)
		}
		fp.Close()
	} else {
		log.Printf("%s-%s has no external port, Skipping Ingress\n", application_name, AppObj.Name)
	}

}
