---
name: Bug report
about: Create a report to help us improve
title: ''
labels: ''
assignees: ''

---

### Prerequisites

* [ ] Are you running the latest version?
* [ ] Have you checked for duplicate bug reports?


### Observed Behaviour
A clear and concise description of what the bug is.

### Steps to reproduce
Steps to reproduce the behavior:
1. Run `./anchorctl test -f ./samples/kube-test.yaml -v 5`
2. With config file:
```yaml
apiVersion: anchor.covarity.dev/alpha1v1
kind: KubeTest
metadata:
  name: opa-validation
...
```

### Expected behavior
A clear and concise description of what you expected to happen.

### Screenshots
If applicable, add screenshots to help explain your problem.

### Context
- OS: [e.g. macOS Mojave ]
- Version [e.g. 22]
- Additional information:
Add any other context about the problem here.
