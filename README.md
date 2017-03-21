# PatientSky Template Renderer ( PSTR )

Templates/


## --all <bool>
Renders Deployment, Service, Autoscaler & Ingress from Json

## --limit <string>
Limits to specific service name, can be used in conjunction with --deploy, --service, --autoscaler & --ingress to narrow output

## --deploy <bool>
Renders only deployment

## --service <bool>
Renders only service

## --autoscaler
Renders only HPA ( AutoScaler )

## --ingress
Renders only Ingress rules

## --namespace <string>
Set namespace to use in template rendering

## --build <string>
Set build. This propagates the value of "Deploy_build"

## --hostname <string>
Comma delimited list of hostnames to use for template rendering.

**WARNING, If multiple services has Service.#.Ports.External.HTTP set this will generate multiple ingress rules with the same hostname if --limit is not used together with**


## notepad
bamboo_deploy_release="34" build_id="f543gr45" bamboo_bulildNumber="4324" CONSUL_APPLICATION="consul_app" cluster_ip="127.0.0.1" CONSUL_ENVIRONMENT="consul_env" CONSUL_PASSWORD="consul_pass" bamboo_buildNumber="123" CONSUL_URL="http://consul" CONSUL_USERNAME="consul_user" git_repo="http://git.repo" ssh_key="rsa1234" NEW_RELIC_LICENSE_KEY="9876er54321" go run *.go --build b1337 --namespace=pltest --autoscaler --hostname=www.roffe.nu,korv.asdf.com --file ./serviceDefinition.json


bamboo_deploy_release="34" build_id="f543gr45" bamboo_bulildNumber="4324" CONSUL_APPLICATION="consul_app" cluster_ip="127.0.0.1" CONSUL_ENVIRONMENT="consul_env" CONSUL_PASSWORD="consul_pass" bamboo_buildNumber="123" CONSUL_URL="http://consul" CONSUL_USERNAME="consul_user" git_repo="http://git.repo" ssh_key="rsa1234" NEW_RELIC_LICENSE_KEY="9876er54321" pstr --build b1337 --namespace=pltest --autoscaler --hostname=www.roffe.nu,korv.asdf.com --file ./serviceDefinition.json