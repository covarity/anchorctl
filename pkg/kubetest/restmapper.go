package kubetest

//
//import (
//	"fmt"
//	"os"
//	"io"
//	"encoding/json"
//	//"time"
//
//	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
//	"k8s.io/client-go/discovery"
//	"k8s.io/client-go/dynamic"
//	"k8s.io/apimachinery/pkg/util/yaml"
//	"k8s.io/apimachinery/pkg/runtime"
//	"k8s.io/apimachinery/pkg/runtime/schema"
//	"k8s.io/apimachinery/pkg/api/meta"
//	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
//
//)
//
//func apply(client *kubernetes.Clientset, file string) {
//
//	f,err := os.Open(file)
//	if err!=nil {
//		//log.Fatal(err)
//	}
//	d := yaml.NewYAMLOrJSONDecoder(f,4096)
//	dd := client.Discovery()
//	apigroups,err := discovery.GetAPIGroupResources(dd)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	restmapper := discovery..NewRESTMapper(apigroups,meta.InterfacesForUnstructured)
//
//
//	for {
//		// https://github.com/kubernetes/apimachinery/blob/master/pkg/runtime/types.go
//		ext := runtime.RawExtension{}
//		if err := d.Decode(&ext); err!=nil {
//			if err == io.EOF {
//				break
//			}
//			log.Fatal(err)
//		}
//		fmt.Println("raw: ",string(ext.Raw))
//		versions := &runtime.VersionedObjects{}
//		//_, gvk, err := objectdecoder.Decode(ext.Raw,nil,versions)
//		obj, gvk, err := unstructured.UnstructuredJSONScheme.Decode(ext.Raw,nil,versions)
//		fmt.Println("obj: ",obj)
//
//		// https://github.com/kubernetes/apimachinery/blob/master/pkg/api/meta/interfaces.go
//		mapping, err := restmapper.RESTMappings(gvk.GroupKind(), gvk.Version)
//		if err != nil {
//			//log.Fatal(err)
//		}
//
//		restconfig := config
//		restconfig.GroupVersion = &schema.GroupVersion {
//			Group: mapping.GroupVersionKind.Group,
//			Version: mapping.GroupVersionKind.Version,
//		}
//		dclient,err := dynamic.NewClient(restconfig)
//		if err != nil {
//			//log.Fatal(err)
//		}
//
//		// https://github.com/kubernetes/client-go/blob/master/discovery/discovery_client.go
//		apiresourcelist, err := dd.ServerResources()
//		if err != nil {
//			//log.Fatal(err)
//		}
//		var myapiresource metav1.APIResource
//		for _,apiresourcegroup := range(apiresourcelist) {
//			if apiresourcegroup.GroupVersion == mapping.GroupVersionKind.Version {
//				for _,apiresource := range(apiresourcegroup.APIResources) {
//					//fmt.Println(apiresource)
//
//					if apiresource.Name == mapping.Resource && apiresource.Kind == mapping.GroupVersionKind.Kind {
//						myapiresource = apiresource
//					}
//				}
//			}
//		}
//		fmt.Println(myapiresource)
//		// https://github.com/kubernetes/client-go/blob/master/dynamic/client.go
//
//		var unstruct unstructured.Unstructured
//		unstruct.Object = make(map[string]interface{})
//		var blob interface{}
//		if err := json.Unmarshal(ext.Raw,&blob); err != nil {
//			//log.Fatal(err)
//		}
//		unstruct.Object = blob.(map[string]interface{})
//		fmt.Println("unstruct:",unstruct)
//		ns := "default"
//		if md,ok := unstruct.Object["metadata"]; ok {
//			metadata := md.(map[string]interface{})
//			if internalns,ok := metadata["namespace"]; ok {
//				ns = internalns.(string)
//			}
//		}
//		res := dclient.Resource(&myapiresource,ns)
//		fmt.Println(res)
//		us,err := res.Create(&unstruct)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println("unstruct response:",us)
//
//
//	}
//}

//https://stackoverflow.com/questions/47116811/client-go-parse-kubernetes-json-files-to-k8s-structures/47139247#47139247

//func parseK8sYaml(fileR []byte) []runtime.Object {
//
//acceptedK8sTypes := regexp.MustCompile(`(Role|ClusterRole|RoleBinding|ClusterRoleBinding|ServiceAccount)`)
//fileAsString := string(fileR[:])
//sepYamlfiles := strings.Split(fileAsString, "---")
//retVal := make([]runtime.Object, 0, len(sepYamlfiles))
//for _, f := range sepYamlfiles {
//if f == "\n" || f == "" {
//// ignore empty cases
//continue
//}
//
//decode := scheme.Codecs.UniversalDeserializer().Decode
//obj, groupVersionKind, err := decode([]byte(f), nil, nil)
//
//if err != nil {
//log.Println(fmt.Sprintf("Error while decoding YAML object. Err was: %s", err))
//continue
//}
//
//if !acceptedK8sTypes.MatchString(groupVersionKind.Kind) {
//log.Printf("The custom-roles configMap contained K8s object types which are not supported! Skipping object with type: %s", groupVersionKind.Kind)
//} else {
//retVal = append(retVal, obj)
//}
//
//}
//return retVal
//}

//https://github.com/kubernetes/client-go/issues/193

// create the dynamic client from kubeconfig
//dynamicClient, err := dynamic.NewForConfig(kubeconfig)
//if err != nil {
//return err
//}
//
//// convert the runtime.Object to unstructured.Unstructured
//unstructuredObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
//if err != nil {
//return err
//}
//
//// create the object using the dynamic client
//nodeResource := schema.GroupVersionResource{Version: "v1", Resource: "Node"}
//createdUnstructuredObj, err := dynamicClient.Resource(nodeResource).Namespace(ns).Create(unstructuredObj)
//if err != nil {
//return err
//}
//
//// convert unstructured.Unstructured to a Node
//var node *corev1.Node
//if err = runtime.DefaultUnstructuredConverter.FromUnstructured(createdUnstructuredObj, node); err != nil {
//return err
//}

//dynClient, err := dynamic.NewForConfig(config)
//...
//clientset, err := kubernetes.NewForConfig(config)
//...
//gvk := obj.GroupVersionKind()
//gk := schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind}
//groupResources, err := restmapper.GetAPIGroupResources(clientset.Discovery())
//...
//rm := restmapper.NewDiscoveryRESTMapper(groupResources)
//mapping, err := rm.RESTMapping(gk, gvk.Version)
//...
//dynClient.Resource(mapping.Resource).Namespace("default").Create(obj, metav1.CreateOptions{})

//https://stackoverflow.com/questions/53341727/how-to-submit-generic-runtime-object-to-kubernetes-api-using-client-go
