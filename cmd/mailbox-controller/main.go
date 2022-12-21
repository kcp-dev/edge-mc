/*
Copyright 2022 The KCP Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/spf13/pflag"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	kcpclient "github.com/kcp-dev/kcp/pkg/client/clientset/versioned"
	kcpinformers "github.com/kcp-dev/kcp/pkg/client/informers/externalversions"

	edgeindexers "github.com/kcp-dev/edge-mc/pkg/indexers"
)

func main() {
	resyncPeriod := time.Duration(0)
	fs := pflag.NewFlagSet("mailbox-controller", pflag.ExitOnError)
	klog.InitFlags(flag.CommandLine)
	fs.AddGoFlagSet(flag.CommandLine)

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	fs.StringVar(&loadingRules.ExplicitPath, "kubeconfig", loadingRules.ExplicitPath, "pathname of kubeconfig file")
	fs.Parse(os.Args[1:])

	ctx := context.Background()
	logger := klog.FromContext(ctx)

	// create cfg
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	cfg, err := kubeConfig.ClientConfig()
	if err != nil {
		logger.Error(err, "failed to make config, if running out of cluster, make sure $KUBECONFIG points to kcp server")
		os.Exit(1)
	}

	cfg.UserAgent = "mailbox-controller"

	// create kcpSharedInformerFactory
	kcpClientset, err := kcpclient.NewForConfig(cfg)
	if err != nil {
		logger.Error(err, "failed to create kcp cluster client")
		os.Exit(1)
	}
	kcpSharedInformerFactory := kcpinformers.NewSharedInformerFactoryWithOptions(
		kcpClientset,
		resyncPeriod,
		kcpinformers.WithExtraClusterScopedIndexers(edgeindexers.ClusterScoped()),
		kcpinformers.WithExtraNamespaceScopedIndexers(edgeindexers.NamespaceScoped()),
	)

	wsInformer := kcpSharedInformerFactory.Tenancy().V1beta1().Workspaces().Informer()
	onAdd := func(obj any) {
		logger.Info("Observed add", "obj", obj)
	}
	onUpdate := func(oldObj, newObj any) {
		logger.Info("Observed update", "oldObj", oldObj, "newObj", newObj)
	}
	onDelete := func(obj any) {
		logger.Info("Observed delete", "obj", obj)
	}
	wsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})
	kcpSharedInformerFactory.Start(ctx.Done())
	<-ctx.Done()
	logger.Info("Time to stop")
}
