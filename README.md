# ECS tools

[![Travis](https://img.shields.io/travis/flou/ecs.svg)](https://travis-ci.org/flou/ecs)
[![Go Report Card](https://goreportcard.com/badge/github.com/flou/ecs)](https://goreportcard.com/report/github.com/flou/ecs)

## Installation

On MacOS:

```
$ curl -LO https://github.com/flou/ecs/releases/download/$(curl -s https://api.github.com/repos/flou/ecs/releases/latest | grep tag_name | cut -d '"' -f 4)/ecs_darwin_amd64
$ chmod +x ecs_darwin_amd64
$ sudo mv ecs_darwin_amd64 /usr/local/bin/ecs
```

On Linux:

```
$ curl -LO https://github.com/flou/ecs/releases/download/$(curl -s https://api.github.com/repos/flou/ecs/releases/latest | grep tag_name | cut -d '"' -f 4)/ecs_linux_amd64
$ chmod +x ecs_linux_amd64
$ sudo mv ecs_linux_amd64 /usr/local/bin/ecs
```

Or, if you have a working Go environment:

```
$ go get -u github.com/flou/ecs
```

### Usage

```
Command line tools to interact with your ECS clusters

Usage:
  ecs [command]

Available Commands:
  help        Help about any command
  image       Print the Docker image of a service running in ECS
  instances   List container instances in your ECS clusters
  services    List services in your ECS clusters
  tasks       List tasks running in your ECS clusters
  update      Update the service to a specific DesiredCount

Flags:
  -h, --help            help for ecs
      --region string   AWS region
      --version         version for ecs

Use "ecs [command] --help" for more information about a command.
```

## List running ECS services

ECS services lists unhealthy services in your ECS clusters

The script will list the clusters and fetch details about services that are
running in it.

Then it determines if the service is healthy or not:

* `[OK]` means that its desiredCount equals its runningCount and it has reached a
  steady state
* `[WARN]` means that its desiredCount and its runningCount equal 0 but it still
  has reached a steady state
* `[KO]` means that its desiredCount does not equal its runningCount, or that it
  has not reached a steady state

### Usage

```
List services in your ECS clusters

Usage:
  ecs services [flags]

Flags:
  -a, --all              Print all services, ignoring their status
  -c, --cluster string   Filter by the name of the ECS cluster
  -h, --help             help for services
  -l, --long             Enable detailed output of containers parameters
  -s, --service string   Filter by the name of the ECS service

Global Flags:
      --region string   AWS region
```

List unhealthy services in all ECS clusters:

```
$ ecs services
--- CLUSTER: ecs-mycluster-dev (listing 4/9 services)
[KO]   tools-hound-dev-1                                    ACTIVE   running 0/1  (hound-dev:75)
[WARN] tools-jenkins-dev-1                                  ACTIVE   running 0/0  (jenkins-dev:247)
[WARN] tools-kibana-dev-1                                   ACTIVE   running 0/0  (srv-kibana:55)
[KO]   tools-sonar-dev-1                                    ACTIVE   running 2/1  (srv-sonar:923)

--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:142)
```


By default `ecs services` only shows services that are `KO` or `WARN`, use the `-a/--all` option to list all services. Also if all services in a cluster are `OK`, all services are shown.

List unhealthy services in a specific ECS cluster:

```
$ ecs services --cluster ecs-mycluster-prod
--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:247)
```

You can also get more information by using the -l/--long option:

```
$ ecs services --long --cluster prod
--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:142)
- Container: jenkins
  Image: 123456789012.dkr.ecr.us-east-1.amazonaws.com/acme/jenkins:2.77-custom
  Memory: 1024 / CPU: 512
  Ports: ->8080
  Environment:
   - JAVA_OPTS: -Dhudson.footerURL=http://mycompany.com
   - JENKINS_OPTS: --prefix=/jenkins
   - JENKINS_SLAVE_AGENT_PORT: 50001
   - PLATFORM: prod
   - PROJECT: acme
```

You can also choose to filter the services listed by their name:

```
$ ecs services -s jenkins
--- CLUSTER: ecs-mycluster-dev (listing 1/9 services)
[WARN] tools-jenkins-dev-1                                  ACTIVE   running 0/0  (jenkins-dev:247)

--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:142)
```

Also, you can combine the flags `-s` and `-c` to filter down a specific service
and check its health across a restricted set of clusters.

## List container instances in ECS clusters

```
List container instances in your ECS clusters

Usage:
  ecs instances [flags]

Flags:
  -c, --cluster string   Filter by the name of the ECS cluster
  -h, --help             help for instances
  -l, --long             Enable detailed output of containers instances

Global Flags:
      --region string   AWS region
```

Example:

```
$ ecs instances -c ecs-mycluster
--- CLUSTER: ecs-mycluster-dev (2 registered instances)
INSTANCE ID           STATUS   TASKS  CPU/used CPU/free  MEM/used MEM/free       PRIVATE IP    INST.TYPE  AGENT  IMAGE         NAME
i-0a2cc6d9443941234   ACTIVE       3       768     1280      1152     2800       10.0.98.85    t2.medium  true   ami-0693ed7f  asg-ecs-mycluster-dev
i-020ff52ddb0538d4a   ACTIVE       3       384     1664       768     3184     10.0.125.146    t2.medium  true   ami-0693ed7f  asg-ecs-mycluster-dev

--- CLUSTER: ecs-mycluster-prod (3 registered instances)
INSTANCE ID           STATUS   TASKS  CPU/used CPU/free  MEM/used MEM/free       PRIVATE IP    INST.TYPE  AGENT  IMAGE         NAME
i-01fad74c0f1b57b85   ACTIVE      11       704     3392      7072     8976      10.0.105.39    m4.xlarge  true   ami-0693ed7f  asg-ecs-mycluster-prod
i-0c75cf9cee1cb5d9e   ACTIVE       6       384     3712      4352    11696      10.0.118.76    m4.xlarge  true   ami-0693ed7f  asg-ecs-mycluster-prod
i-0c61781827ef44a52   ACTIVE       2       128     3968      1536    14512     10.0.104.249    m4.xlarge  true   ami-0693ed7f  asg-ecs-mycluster-prod
```

## Update an ECS service

```
Update the service to a specific DesiredCount

Usage:
  ecs update [flags]

Flags:
  -c, --cluster string   Name of the ECS cluster
      --count int        New DesiredCount (default -1)
  -f, --force            Force a new deployment of the service
  -h, --help             help for update
  -s, --service string   Name of the ECS service

Global Flags:
      --region string   AWS region
```

Example:

```
ecs update --cluster ecs-mycluster-prod --service tools-jenkins-prod-1 --count 0
```

## List tasks running on ECS

```
List tasks running in your ECS clusters

Usage:
  ecs tasks [flags]

Flags:
  -c, --cluster string   Filter by the name of the ECS cluster
  -h, --help             help for tasks
  -l, --long             Enable detailed output of containers parameters
  -s, --service string   Filter by the name of the ECS service

Global Flags:
      --region string   AWS region
```


## Find the images in an ECS service

```
Print the Docker image of a service running in ECS

Usage:
  ecs image [flags]

Flags:
      --cluster string   Name of the ECS cluster
  -h, --help             help for image
      --service string   Name of the ECS service

Global Flags:
      --region string   AWS region
```
