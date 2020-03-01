package assertexec

import (
	"anchorctl/pkg/kubetest/common"
	"anchorctl/pkg/kubetest/resource"
	"anchorctl/pkg/logging"
	"anchorctl/pkg/resultaggregator"
	"bytes"
	"fmt"
	"github.com/mitchellh/mapstructure"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
)

type execTest struct {
	Command  []string `yaml:"command"`
	Contains string
	client   *kubernetes.Clientset
}

var log *logging.Logger

func Parse(clientset *kubernetes.Clientset, logger *logging.Logger, input interface{}) (common.KubeTester, error) {
	var execTest *execTest
	log = logger
	err := mapstructure.Decode(input, &execTest)
	if err != nil {
		return nil, err
	}
	execTest.client = clientset
	return execTest, nil
}

func (et *execTest) Valid(res *resource.Resource) (bool, error) {

	if len(et.Command) < 1 {
		return false, fmt.Errorf("length of command less than 1")
	}

	return common.ValidTest("JSONTest", map[string]string{
		"JSONPath":                 et.Contains,
		"ObjectRef.Spec.Kind":      res.ObjectRef.Spec.Kind,
		"ObjectRef.Spec.Namespace": res.ObjectRef.Spec.Namespace,
	})
}

func (et *execTest) Test(res *resource.Resource) *resultaggregator.TestRun {
	testRun := resultaggregator.NewTestRun("AssertExec", strings.Join(et.Command, " ")+" returns "+et.Contains)

	output, err := execCommand(et.client, res, et.Command)
	if err != nil {
		testRun.Invalid = true
		testRun.Details = err.Error()
		return testRun
	}

	if strings.Contains(output, et.Contains) {
		testRun.Passed = true
		log.InfoWithFields(map[string]interface{}{
			"Test":    "AssertExec",
			"Command": strings.Join(et.Command, " "),
			"result":  et.Contains,
			"status":  "PASSED",
		}, "AssertExec return the expected output.")
	}

	return testRun
}

// ExecCommands executes arbitrary commands inside the given container.
// If developers set the containerName = "", this function will choose
// the first container which is listed from given pod as the default container to run the commands
// similar to kubectl command 'k exec -it <pod.Name> -n <pod.Namespace> -- <commands>'
func execCommand(client *kubernetes.Clientset, res *resource.Resource, commands []string) (string, error) {

	var execOut, execErr bytes.Buffer
	targetContainerIndex := -1

	objRef := res.ObjectRef

	pods, _ := client.CoreV1().Pods(objRef.Spec.Namespace).List(metav1.ListOptions{
		LabelSelector: common.MapToString(objRef.Spec.Labels, ","),
	})

	if len(pods.Items) < 1 {
		err := fmt.Errorf("No Pods founds in namespace " + objRef.Spec.Namespace + "with labels " + common.MapToString(objRef.Spec.Labels, " "))
		log.Error(err, "Assert Exec failed with error")
		return "", err
	}

	pod := pods.Items[0]

	// If not provide containerName, will choose the first one as default.
	if objRef.Spec.Container == "" {
		targetContainerIndex = 0
	} else {
		// Iterate through all containers looking for the one with containerName.
		for i, cr := range pod.Spec.Containers {
			if cr.Name == objRef.Spec.Container {
				targetContainerIndex = i
				break
			}
		}

		if targetContainerIndex < 0 {
			err := fmt.Errorf("could not find %s container to exec to", objRef.Spec.Container)
			log.Error(err, "Assert Exec failed with error")
			return "", err
		}
	}

	req := client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(pod.Namespace).
		SubResource("exec")
	req.VersionedParams(&v1.PodExecOptions{
		Container: pod.Spec.Containers[targetContainerIndex].Name,
		Command:   commands,
		Stdout:    true,
		Stderr:    true,
	}, scheme.ParameterCodec)

	exec, _ := remotecommand.NewSPDYExecutor(common.RestKubeConfig, "POST", req.URL())

	err := exec.Stream(remotecommand.StreamOptions{
		Stdout: &execOut,
		Stderr: &execErr,
		Tty:    false,
	})

	return execOut.String(), err
}
