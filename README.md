# ECS tools

[![Go Report Card](https://goreportcard.com/badge/github.com/flou/ecs)](https://goreportcard.com/report/github.com/flou/ecs)

## Installation

```
go get -u github.com/flou/ecs
```

### Usage

```
usage: ecs [<flags>] <command> [<args> ...]

ECS Tools

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -r, --region=REGION  AWS Region

Commands:
  help [<command>...]
    Show help.

  monitor [<flags>]
    List unhealthy services in your ECS clusters

  scale --cluster=CLUSTER --service=SERVICE --count=COUNT
    Scale the service to a specific DesiredCount

  image --cluster=CLUSTER --service=SERVICE
    Return the Docker image of a service running in ECS
```

## Monitor

ECS monitor lists unhealthy services in your ECS clusters

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
usage: ecs monitor [<flags>]

List unhealthy services in your ECS clusters

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -r, --region=REGION    AWS Region
      --cluster=CLUSTER  Select the ECS cluster to monitor
  -f, --filter=FILTER    Filter by the name of the ECS cluster
  -l, --long             Enable detailed output of containers parameters
  -a, --all              Enable detailed output of containers parameters
```

List unhealthy services in all ECS clusters:

```
$ ecs monitor
--- CLUSTER: ecs-mycluster-dev (listing 4/9 services)
[KO]   tools-hound-dev-1                                   ACTIVE   running 0/1  (hound-dev:75)
[WARN] tools-jenkins-dev-1                                 ACTIVE   running 0/0  (jenkins-dev:247)
[WARN] tools-kibana-dev-1                                  ACTIVE   running 0/0  (srv-kibana:55)
[KO]   tools-sonar-dev-1                                   ACTIVE   running 2/1  (srv-sonar:923)

--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:142)
```

By default `ecs monitor` only shows services that are `KO` or `WARN`, use the `-a/--all` option to list all services. Also if all services in a cluster are `OK`, all services are shown.

List unhealthy services in a specific ECS cluster:

```
$ ecs monitor --cluster ecs-mycluster-prod
--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:247)
```

You can also get more information by using the -l/--long option:

```
$ ecs monitor --long --filter prod
--- CLUSTER: ecs-mycluster-prod (listing 1/12 services)
[WARN] tools-jenkins-prod-1                                 ACTIVE   running 0/0  (jenkins-prod:142
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

## Scale
