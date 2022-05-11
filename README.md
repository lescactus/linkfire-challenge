# Linkfire Challenge [![k8s](https://github.com/lescactus/linkfire-challenge/actions/workflows/k8s.yml/badge.svg)](https://github.com/lescactus/linkfire-challenge/actions/workflows/k8s.yml) [![Docker Image CI](https://github.com/lescactus/linkfire-challenge/actions/workflows/docker-image.yml/badge.svg)](https://github.com/lescactus/linkfire-challenge/actions/workflows/docker-image.yml) [![Go](https://github.com/lescactus/linkfire-challenge/actions/workflows/go.yml/badge.svg)](https://github.com/lescactus/linkfire-challenge/actions/workflows/go.yml)

## Instructions

* Create or use an existent application, it must be dockerized.
* Create the necessary Kubernetes manifests
* Create a pipeline to build and deploy this app into Kubernetes

## Solution

`linkfire-challenge` is a very minimalistic web app written in [Golang](https://golang.org/).
The repository has the following structure:

```
.                       // Root of the project. Contains source code, Dockerfile, licence, etc ...
├── *.go
├── .github/            // Contains GitHub CI/CD worklows
├── deploy/             // Contains Kubernetes manifests
│   └── *.yaml
├── Dockerfile
├── LICENSE
├── README.md
├── skaffold.yaml
└── ...
```

### Specifications

`GET /rest/ready` returns whether the app is ready to receive requests

`GET /rest/alive` returns whether the app is alive

`GET /rest/v1/ping` returns a pong and some information about the server

`POST /rest/v1/hello` returns a custom and some information about the server

**Example**:

```
> GET /rest/ready HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/health+json
< Date: Wed, 11 May 2022 14:42:24 GMT
< Content-Length: 17
< 
{"status":"pass"}
```

```
> GET /rest/alive HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/health+json
< Date: Wed, 11 May 2022 14:42:24 GMT
< Content-Length: 17
< 
{"status":"pass"}
```

```
> GET /rest/v1/ping HTTP/1.1

< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 11 May 2022 14:43:21 GMT
< Content-Length: 168
< 
{"hostname":"linkfire-challenge-6b69c9c577-49sr9","message":"Pong","goos":"linux","goarch":"amd64","runtime":"go1.16.7","cpu":8,"application_name":"linkfire-challenge"}     
```

```
> POST /rest/v1/hello HTTP/1.1
> {"message":"Hello Linkfire!"}

< HTTP/1.1 200 OK
< Content-Type: application/json
< Date: Wed, 11 May 2022 14:44:13 GMT
< Content-Length: 179
<
{"hostname":"linkfire-challenge-6b69c9c577-49sr9","message":"Hello Linkfire!","goos":"linux","goarch":"amd64","runtime":"go1.16.7","cpu":8,"application_name":"linkfire-challenge"}    
```

### Configuration

`linkfire-challenge` is a 12-factor app using [Viper](https://github.com/spf13/viper) as a configuration manager. It can read configuration from environment variables. Reading from a configuration file could be implemented in the future.

#### Available variables

* `APP_ADDR`(default value: `:8080`). Define the TCP address for the server to listen on, in the form "host:port".

* `SERVER_READ_TIMEOUT` (default value: `30s`). Maximum duration for reading the entire request, including the body (`ReadTimeout`).

* `SERVER_READ_HEADER_TIMEOUT` (default value: `10s`). Amount of time allowed to read request headers (`ReadHeaderTimeout`).

* `SERVER_WRITE_TIMEOUT` (default value: `30s`). Maximum duration before timing out writes of the response (`WriteTimeout`).


### Building `linkfire-challenge`

<details>
<summary>Click to expand</summary>

#### From source with go

You need a working [go](https://golang.org/doc/install) toolchain (It has been developped and tested with go 1.16 and go 1.16 only, but should work with go >= 1.14). Refer to the official documentation for more information (or from your Linux/Mac/Windows distribution documentation to install it from your favorite package manager).

```sh
# Build from sources. Use the '-o' flag to change the compiled binary name
go build

# Default compiled binary is linkfire-challenge
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./linkfire-challenge
```

#### From source with docker

If you don't have [go](https://golang.org/) installed but have docker, run the following command to build inside a docker container:

```sh
# Build from sources inside a docker container. Use the '-o' flag to change the compiled binary name
# Warning: the compiled binary belongs to root:root
docker run --rm -it -v "$PWD":/app -w /app golang:1.16 go build

# Default compiled binary is linkfire-challenge
# You can optionnaly move it somewhere in your $PATH to access it shell wide
./linkfire-challenge
```

#### From source with docker but built inside a docker image

If you don't want to pollute your computer with another program, `linkfire-challenge` comes with its own docker image:

```sh
docker build -t linkfire-challenge .

docker run --rm -p 8080:8080 linkfire-challenge
```

#### Unit tests

To run the test suite, run the following commands:

```sh
# Run the unit tests. Remove the '-v' flag to reduce verbosity
go test -v ./... 

# Get coverage to html format
go test -coverprofile -v /tmp/cover.out ./...
go tool cover -html=/tmp/cover.out -o /tmp/cover.out.html
```

</details>

### Kubernetes deployment

#### With Skaffold

Ensure you have a properly working and accessible Kubernetes cluster with a valid `~/.kube/config`. 
This project is using [Skaffold](https://skaffold.dev/) to deploy to a local Kubernetes cluster, such as Minikube or KinD. You can dowload `skaffold` [here](https://skaffold.dev/docs/install/#standalone-binary). It is assumed that `skaffold` is installed.

To deploy to a local Kubernetes cluster, simply run `skaffold run`:

<details>
<summary>Click to expand</summary>
```
$ skaffold run 
Generating tags...
 - linkfire-challenge -> linkfire-challenge:62fb7fa-dirty
Checking cache...
 - linkfire-challenge: Not found. Building
Starting build...
Found [kind-kind] context, using local docker daemon.
Building [linkfire-challenge]...
Target platforms: [linux/amd64]
[+] Building 23.3s (16/16) FINISHED                                                                                                                                                                                                            
 => [internal] load build definition from Dockerfile                                                                                                                                                                                      0.0s
 => => transferring dockerfile: 38B                                                                                                                                                                                                       0.0s
 => [internal] load .dockerignore                                                                                                                                                                                                         0.0s
 => => transferring context: 34B                                                                                                                                                                                                          0.0s
 => [internal] load metadata for docker.io/library/alpine:3                                                                                                                                                                               0.0s
 => [internal] load metadata for docker.io/library/golang:1.16-alpine                                                                                                                                                                     0.0s
 => [builder 1/6] FROM docker.io/library/golang:1.16-alpine                                                                                                                                                                               0.0s
 => [internal] load build context                                                                                                                                                                                                         0.0s
 => => transferring context: 12.99kB                                                                                                                                                                                                      0.0s
 => [stage-1 1/4] FROM docker.io/library/alpine:3                                                                                                                                                                                         0.0s
 => CACHED [builder 2/6] ADD go.* /go/src/                                                                                                                                                                                                0.0s
 => CACHED [builder 3/6] WORKDIR /go/src/                                                                                                                                                                                                 0.0s
 => [builder 4/6] RUN go mod download                                                                                                                                                                                                    19.9s
 => [builder 5/6] COPY . /go/src/                                                                                                                                                                                                         0.2s
 => [builder 6/6] RUN go build -o main                                                                                                                                                                                                    3.0s
 => CACHED [stage-1 2/4] RUN apk update     && apk add ca-certificates     && rm -rf /var/cache/apk*     && adduser -u 1000 -D -s /bin/sh app     && install -d -m 0750 -o app -g app /app                                                0.0s
 => CACHED [stage-1 3/4] WORKDIR /app                                                                                                                                                                                                     0.0s
 => CACHED [stage-1 4/4] COPY --from=builder /go/src/main /app                                                                                                                                                                            0.0s
 => exporting to image                                                                                                                                                                                                                    0.0s
 => => exporting layers                                                                                                                                                                                                                   0.0s
 => => writing image sha256:6276a5062a99b7ae078ba5dbddb108029258fd93a887a53332ed31fc89e318c3                                                                                                                                              0.0s
 => => naming to docker.io/library/linkfire-challenge:62fb7fa-dirty                                                                                                                                                                       0.0s
Build [linkfire-challenge] succeeded
Starting test...
Tags used in deployment:
 - linkfire-challenge -> linkfire-challenge:6276a5062a99b7ae078ba5dbddb108029258fd93a887a53332ed31fc89e318c3
Starting deploy...
Loading images into kind cluster nodes...
 - linkfire-challenge:6276a5062a99b7ae078ba5dbddb108029258fd93a887a53332ed31fc89e318c3 -> Loaded
Images loaded in 1.496 second
 - deployment.apps/linkfire-challenge created
 - service/linkfire-challenge created
 - serviceaccount/linkfire-challenge created
Waiting for deployments to stabilize...
 - deployment/linkfire-challenge: waiting for rollout to finish: 0 of 1 updated replicas are available...
 - deployment/linkfire-challenge is ready.
Deployments stabilized in 7.092 seconds
You can also run [skaffold run --tail] to get the logs
```

</details>


The following happened:

* Skaffold will generate a docker tag based on the current commit (This is [customizable](https://skaffold.dev/docs/references/yaml/#build-tagPolicy)).

* If the image doesn't exist locally, Skaffold will build it.

* Once the docker image is built, Skaffold will substitute the raw image defined in `deploy/deployment.yaml` with the image just built.

* Skaffold will apply the manifests in `deploy/`.

* Skaffold will wait for the `linkfire-challenge` deployment to be ready.

For more informations about Skaffold and what it can do, visit the project [documentation](https://skaffold.dev/docs/).

#### Without Skaffold

To deploy to a Kubernetes cluster without Skaffold, simply build & push the docker image to an external registry. Then change the docker image name to include the registry in the `deploy/deployment.yaml` manifest.

#### Accessing `linkfire-challenge` REST API

Since the service associated with `linkfire-challenge` is of type [`LoadBalancer`](https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer), it is possible to reach it using its external ingress IP. To curl it, run the following command:

```sh
# Typically on AWS
hostname="$(kubectl get svc -n antaeus antaeus -ojsonpath='{.status.loadBalancer.ingress[0].hostname}')"
curl "http://${hostname}:80/rest/xxx"


# Typically on Minikube or GCP
ip="$(kubectl get svc -n antaeus antaeus -ojsonpath='{.status.loadBalancer.ingress[0].ip}')"
curl "http://${ip}:80/rest/xxx"
```

### CI/CD

This project is using GitHub Actions as a CI/CD. The workflows are decclared in `.github/workflows/`

* `docker.yaml`: Simply build the Docker image using the Dockerfile. This job can be extended if needed to also push it to an external registry.

* `go.yaml`: Compile the app with `go build`, run the unit tests and ensure no race condition is detected. Typically, this job can be extended with format vet, linter, etc ...

* `k8s.yaml`: This job will:
    * Create a simple Kubernetes cluster with [`KinD`](https://kind.sigs.k8s.io/)
    * Install [MetalLB](https://metallb.org/) in the cluster to be able to provision `LoadBalancer` services
    * Install Skaffold
    * Execute the `skaffold run` command to deploy `linkfire-challenge` in KinD
    * Run some very basic e2e tests with curl:
        * `GET /rest/ready`
        * `GET /rest/alive`
        * `GET /rest/v1/ping`
        * `POST /rest/v1/hello`

### Improvements area

<details>
<summary>Click to expand</summary>

`linkfire-challenge` is a very minimalistic service far from being "production ready" (and it could even be written with less lines of code). Some area of improvements might include but not limited to:

* Avoid using `io.ReadAll` in `Hello()` as it loads all the request body into memory, leading to a possible memory exhaustion (and possible Denial Of Service). Better use a `bytes.Buffer` or `io.Copy` instead

* Provide structured (json) logging to stdout with tiered log levels (fatal, error, warn, info, debug, trace) and without sensitive information

* Provide metrics

* Provide APM using a tracing library such as OpenTelemetry

* Provide alerts and monitoring (Grafana) dashboards based on above-mentioned metrics and APM

* Provide incident response checklist

* Provide accurate cpu and memory requests/limits based on stress-test benchmarks

* Use of a [PodDistruptionBudget](https://kubernetes.io/docs/tasks/run-application/configure-pdb/) and scale to multiple replicas with the use of anti-affinity rules to spread accross multiple availability zones

* Be [12factor](https://12factor.net/) compliant by reading configuration at runtime from config maps or secrets (env variables or config files) - such as "log level", "tcp port", etc ...

* Provide a Swagger endpoint

* Support graceful shutdowns for interrupt signals (SIGTERM)

* Ensure the service is stateless by using an external store provider (SQL, blob, NoSQL, k/v, etc...)

* Authenticate API calls (authentication/authorization) or even better, do it either at the edge (via an API gateway for instance) or in a sidecar proxy

* Do not use the `latest` docker tag ([never](https://stevelasker.blog/2018/03/01/docker-tagging-best-practices-for-tagging-and-versioning-docker-images/)). Instead, provide [semantic versioning](https://semver.org/)

* Implement retries policies or circuit-breaker (could ideally also by done by a Service Mesh)

* Improve the CD pipeline to automates the deployment of the new image in a real Kubernetes cluster and run a test suite: integration tests, e2e tests, non-regression tests, stress tests etc ..., with automated logs and metrics analysis to detect anormalities and do environment promotion. [Flagger](https://flagger.app/) is good for that, [Harness](https://harness.io/) too

* To deploy in Kubernetes for production, a [Helm chart](https://helm.sh/) or a [Kustomization](https://kubectl.docs.kubernetes.io/references/kustomize/kustomization/) would come handy, especially when managing multiple environments

* Follow the GitOps principle with tools such as the amazing [FluxCD](https://fluxcd.io/) or [ArgoCD](https://argoproj.github.io/)

* The usage of a `LoadBalancer` service for a single deployment should be discouraged in a production environment. Instead, an [`Ingress`](https://kubernetes.io/docs/concepts/services-networking/ingress/) would be better. Ideally, at the edge should stand an API Gateway or cloud Load Balancer doing TLS termination and redirection, authn/authz (could also be done via sidecar proxy), WAF, audit logs, request validation, etc ...

</details>