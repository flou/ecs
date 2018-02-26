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

## List instances in ECS clusters

```
usage: ecs instances [<flags>]

List container instances in your ECS clusters

Flags:
      --help           Show context-sensitive help (also try --help-long and --help-man).
  -r, --region=REGION  AWS Region
  -f, --filter=FILTER  Filter by the name of the ECS cluster
```

Example:

```
$ ecs instances -f ecs-mycluster
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

## Scale

```
usage: ecs scale --cluster=CLUSTER --service=SERVICE --count=COUNT

Scale the service to a specific DesiredCount

Flags:
      --help             Show context-sensitive help (also try --help-long and --help-man).
  -r, --region=REGION    AWS Region
      --cluster=CLUSTER  Name of the ECS cluster
      --service=SERVICE  Name of the service
      --count=COUNT      New DesiredCount
```

Example:

```
ecs scale --cluster ecs-mycluster-prod --service tools-jenkins-prod-1 --count 0
```
