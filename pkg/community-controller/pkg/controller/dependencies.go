package controller

import (
	"encoding/json"
	"fmt"

	neptuneplus "github.com/lterrac/edge-autoscaler/pkg/apis/neptuneplus/v1alpha1"
	openfaasv1 "github.com/openfaas/faas-netes/pkg/apis/openfaas/v1"
	corev1 "k8s.io/api/core/v1"
)

func (c *CommunityController) GetFunctionDependency(function *openfaasv1.Function) (*corev1.EnvVar, error) {
	fName := function.Name
	fNamespace := function.Namespace

	envVar := &corev1.EnvVar{
		Name:  "DEPDAG_NODE",
		Value: "",
	}

	dag, err := c.listers.DependencyGraphs(fNamespace).Get(fName)
	if err != nil {
		return envVar, err
	}

	var fNode *neptuneplus.FunctionNode = nil
	for _, node := range dag.Spec.Nodes {
		if node.FunctionName == fName && node.FunctionNamespace == fNamespace {
			fNode = node.DeepCopy()
		}
	}
	if fNode == nil {
		return envVar, fmt.Errorf("dag %s/%s does not define dependencies for function %s/%s", dag.Namespace, dag.Name, fNamespace, fName)
	}

	marshaled, err := json.Marshal(*fNode)
	if err != nil {
		return envVar, fmt.Errorf("could not marshal dependency string for function %s/%s", fNamespace, fName)
	}

	envVar.Value = string(marshaled)
	return envVar, nil
}
