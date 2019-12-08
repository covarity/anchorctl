[![codecov](https://codecov.io/gh/trussio/anchorctl/branch/master/graph/badge.svg)](https://codecov.io/gh/trussio/anchorctl)
[![Build Status](https://travis-ci.org/trussio/anchorctl.svg?branch=master)](https://travis-ci.org/trussio/anchorctl)


# Anchorctl

Anchorctl is a command line utility to enable a test driven approach to building distributed systems like Kubernetes. The utility
works in conjunction with Anchor kubernetes controller (Under development) to create a CRD based user experience for
testing systems.

The kinds of tests supported are:
- KubeTest: A collection of tests that are specifically related to the Kubernetes Ecosystem such as asserting the 
functionality of admission controllers, asserting the value of jsonpaths and networkpolicies (under development). 


## Installation

### Binary

The following make entry point build the application and produces a binary.

```bash
make run
```

### Docker

The following make entry point builds a docker image of `anchorctl`. The latest CI build version of the utility can 
be taken from `docker.pkg.github.com/trussio/anchorctl/anchorctl`

```bash
make docker
```

## Usage

### KubeTest

KubeTest contains common function to enable testing common kubernetes features.


The types of Kubernetes tests include:
- `AssertJSONPath`: Takes a json path and a value and assert the json path of the object in cluster is equal to the value. 
Using this type of test, we can test the status of a deployment / pod, the number of replicas and anything else that is 
accessible in the yaml output of Kubernetes objects.

```yaml
- type: AssertJSONPath
  jsonPath: ".spec.nodeName"
  value: "docker-desktop"
```

- `AssertValidation`: Used to ensure that the validation admission controller throw the expected error. Take a file and 
an action, applies the action to the file and assert that the error equals the expected error.

```yaml
- type: AssertValidation
  expectedError: "Internal error occurred: admission webhook \"webhook.openpolicyagent.org\" denied the request: External Loadbalancers cannot be deployed in this cluster"
  objectRef:
    action: "CREATE"
    filePath: "./samples/fixtures/loadbalancer.yaml"
```

- `AssertMutation`: Used to ensure that the mutating admission controller mutates the kubernetes object upon creation 
as expected. 

```yaml
- type: AssertMutation
  jsonPath: ".metadata.labels.function"
  value: "workload"
  objectRef:
    action: "CREATE"
    filePath: "./samples/fixtures/deploy.yaml"
```

### Prerequisites

- go >= 1.13
- Kubernetes Cluster

### Examples


## Why?

With the growing adoption of Kubernetes and increase in the number of varied services that runs inside, the issue of 
platform validation and testing starts to emerge. With no easy mechanism to test, teams face the challenge of how to 
gain confidence with platform changes in an automated fashion. 

We saw the necessity to build a solution to address the lack of distributed systems and application testing. 

The vision for `anchorctl`:

- Provide stability and reliability to the platform
- Help with root cause analysis: first point of reference in case of failure
- Understanding and validating network and service mesh topologies

## Reference 
- Test OPA:
    - [Anchor Test File](https://github.com/covarity/examples/blob/master/examples/test-admission-controller/test/anchor_test_ci.yaml)
    - [Code](https://github.com/covarity/examples/blob/master/examples/test-admission-controller)
- Kube Forum Sydney 2019 Talk: 
    - [Slides](https://github.com/covarity/demos/tree/master/kube-forum-2019)
    - Talk (Coming Soon)
    