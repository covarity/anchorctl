<!-- [![codecov](https://codecov.io/gh/trussio/anchorctl/branch/master/graph/badge.svg)](https://codecov.io/gh/trussio/anchorctl) -->
[![Go Report Card](https://goreportcard.com/badge/github.com/covarity/anchorctl)](https://goreportcard.com/report/github.com/covarity/anchorctl)


# Anchorctl

Anchorctl is a command line utility to enable a test driven approach to developing on distributed systems like Kubernetes. The utility
works in conjunction with Anchor kubernetes controller (Under development) to create a CRD based user experience for testing systems.

The kinds of tests supported are:
- KubeTest: A collection of tests that are specifically related to the Kubernetes Ecosystem such as asserting the functionality of admission controllers, asserting the value of jsonpaths and networkpolicies (under development).


## Installation

### Binary

The following make entry point build the application and produces a binary.

```bash
make run
```

### Docker

The following make entry point builds a docker image of `anchorctl`. The latest CI build version of the utility can be taken from `docker.pkg.github.com/trussio/anchorctl/anchorctl`

```bash
make docker
```

## Resources

Resources provide common mechanisms to refer to an object or a file.

### Manifests

Manifests contain the following fields:
- Path: Relative path to the file with kube resources
- Action: "CREATE", "UPDATE" or "DELETE" action to apply to the file

### ObjectRefs

ObjectRefs provide an interface to communicate with existing objects in the cluster. ObjectRef contains the following fields:
```yaml
type: "Resource"
spec:
    kind: Pod # Kind of kubernetes resource to look for
    namespace: default # Namespace of the resource
    labels:
      hello: world # Label value of the resource
```

With the above information, `anchorctl` is able to find the object from a cluster.

### Lifecycle

Similar to the Pod lifecycle, `anchorctl` support PostStart and PreStop hooks. These hooks can be used to set up your test environment such as create a testing namespace, etc.

#### PostStart && PreStop

PostStart is actioned before the tests are executed. PreStop is actioned after tests are executed and before program exits. The files are actioned on in the order they are structured.

Take a list of Path and Action, such as:
```yaml
lifecycle:
  # Runs before the tests
  postStart:
  - path: "./samples/fixtures/applications-ns.yaml"
    action: "CREATE"
  - path: "./samples/fixtures/hello-world-nginx.yaml"
    action: "CREATE"
  # Runs after the tests
  preStop:
  - path: "./samples/fixtures/hello-world-nginx.yaml"
    action: "DELETE"
  - path: "./samples/fixtures/applications-ns.yaml"
    action: "DELETE"
```

---

## Tests

### KubeTest

KubeTest contains common functions to enable testing common kubernetes features.

The types of Kubernetes tests include:
- `AssertJSONPath`: Takes a jsonpath and a value and asserts that the jsonpath of the objectRef in cluster is equal to the value.
Using this type of test, we can test the status of a deployment / pod, the number of replicas and anything else that is accessible in the yaml output of Kubernetes objects.

```yaml
# Assert that Pods in namespace applications with label hello=world is scheduled on the docker-desktop node.
- type: AssertJSONPath
  spec:
    jsonPath: ".spec.containers[0].image"
    value: "nginx"
  resource:
    objectRef:
      type: Resource
      spec:
        kind: Pod
        namespace: applications
        labels:
          hello: world
```

- `AssertValidation`: Used to ensure that the validation admission controller throw the expected error. Take a file and
an action, applies the action to the file and assert that the error equals the expected error.

```yaml
# Assert that when attempting to create the resources in .resource.manifest.path, the error is returned by the API Server.
- type: AssertValidation
  spec:
    containsResponse: "External Loadbalancers cannot be deployed in this cluster"
  resource:
    manifest:
      path: "./samples/fixtures/loadbalancer.yaml"
      action: CREATE
```

- `AssertMutation`: Used to ensure that the mutating admission controller mutates the kubernetes object upon creation
as expected.

```yaml
# Assert that after creating the resources in .resource.manifest.path, the jsonpath of the object created has defined value.
- type: AssertMutation
  spec:
    jsonPath: ".metadata.labels.function"
    value: "workload"
  resource:
    manifest:
      path: "./samples/fixtures/deploy.yaml"
      action: CREATE

```

### Prerequisites

- go >= 1.13
- Kubernetes Cluster
- kubectl

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
    - [Talk](https://youtu.be/tGDAuij5RvE)
