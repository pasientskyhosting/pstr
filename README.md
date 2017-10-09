# PatientSky Template Renderer ( PSTR )
This project is to render templates used for Kubernetes deployment

## --all <bool>
Renders Deployment, Service, Autoscaler & Ingress from Json

## --limit <string>
Limits to specific service name, can be used in conjunction with --deploy, --service, --autoscaler & --ingress to narrow output

## --deploy <bool>
Renders only deployment

## --service <bool>
Renders only service

## --genericservice <bool>
Renders only generic service

## --autoscaler
Renders only HPA ( AutoScaler )

## --ingress
Renders only Ingress rules

## --cronjob
Renders only Cronjob rules

## --namespace <string>
Set namespace to use in template rendering

## --build <string>
Set build. This propagates the value of "Deploy_build"

## --hostname <string>
Comma delimited list of hostnames to use for template rendering.

## --output <string>
Path to where to write output YAML files.

If not specified files will be written in the current folder


**WARNING, If multiple services has Service.#.Ports.External.HTTP set this will generate multiple ingress rules with the same hostname if --limit is not used**


## notepad
bamboo_CONSUL_APPLICATION=consul_app1 \
bamboo_CONSUL_ENVIRONMENT=consul_environment1 \
bamboo_CONSUL_PASSWORD=consul_password \
bamboo_CONSUL_URL=https://consul-host.ps.com \
bamboo_CONSUL_USERNAME=consul_username \
bamboo_NEW_RELIC_API_KEY_PASSWORD=NR_password \
bamboo_NEW_RELIC_API_URL=https://newrelic.api.url/ \
bamboo_NEW_RELIC_LICENSE_KEY_PASSWORD=NR_license \
bamboo_buildNumber=B12345 \
bamboo_deploy_release=Bamboorelease1 \
cluster_ip=127.0.0.1 \
git_repo=https://git.com/psdev/arepo.git \
ssh_key=fsdfosdhfjdsakhfkljsahfklhsa \
go run *.go  --build_id dfb1337 --namespace=hptest --all --hostname=test.domain.com,test2.another.com --file ./serviceDefinition.json --output ./out
