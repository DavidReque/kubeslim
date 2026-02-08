# kube-slim

_kube-slim is a Kubernetes client library for Go optimized for runtime performance and final binary size_

## Introduction

While working on the [`kubetail`](https://github.com/kubetail-org/kubetail) CLI tool we noticed that our final binary size and that of other tools in the K8s ecosystem (e.g. `kubectl`, `helm`) was large compared to other unix utilities and we traced the problem to the [Kubernetes client-go library](https://github.com/amorey/size-matters) which typically adds 20MB+ to final binaries. We also noticed other issues with `client-go` such as large memory consumption with informers. This library is an attempt to fix those issues by starting fresh and designing a library that's optimized for speed, memory consumption and binary size.

Currently here's the typical binary size for a simple Golang executable with `client-go.Clientset` compared to `kube-slim`:

* With `Clientset` - 38 MB
* With `DynamicClient` - 16 MB
* With `kube-slim` - 14 MB

This library is very much a work in progress and we welcome your help! 

## Installation

```console
go get github.com/kubetail-org/kube-slim
```

## Basic usage

```go
import (
	"context"
	"fmt"
	"log"
	"path/filepath"

  slim "github.com/kubetail-org/kube-slim"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type namespaceList struct {
  Items []struct {
    Metadata struct {
      Name string `json:"name"`
    } `json:"metadata"`
  } `json:"items"`
}

type podList struct {
  Items []struct {
    Metadata struct {
      Name string `json:"name"`
    } `json:"metadata"`
  } `json:"items"`
}

func main() {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")

	// Load configuration from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("failed to load kubeconfig: %v", err)
	}

	// Create the kube-slim client
	client, err := slim.NewForConfig(config)
	if err != nil {
		log.Fatalf("failed to create kube-slim client: %v", err)
	}

  // Get namespaces
  namespaceGVR := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "namespaces"}

	namespaces, err := slim.List[namespaceList](context.TODO(), client, namespaceGVR)
	if err != nil {
		log.Fatalf("failed to list namespaces: %v", err)
	}

	if len(namespaces.Items) == 0 {
		fmt.Println("No namespaces found.")
		return
	}

	for _, ns := range namespaces.Items {
		fmt.Println(ns.Metadata.Name)
	}

  // Get pods
  podGVR := schema.GroupVersionResource{Group: "core", Version: "v1", Resource: "pods"}

	pods, err := slim.List[podList](context.TODO(), client, podGVR)
	if err != nil {
		log.Fatalf("failed to list pods: %v", err)
	}

	if len(pods.Items) == 0 {
		fmt.Println("No pods found.")
		return
	}

	for _, pod := range pods.Items {
		fmt.Println(pod.Metadata.Name)
	}
}
```

## Get Involved

This library is very much a work in progress! If you have ideas on how to improve it or if you'd like help using it you can:

* Create a [GitHub Issue](https://github.com/kubetail-org/kube-slim/issues)
* Send us an email ([hello@kubetail.com](hello@kubetail.com))
* Join our [Discord server](https://discord.gg/CmsmWAVkvX) or [Slack channel](https://join.slack.com/t/kubetail/shared_invite/zt-2cq01cbm8-e1kbLT3EmcLPpHSeoFYm1w)
