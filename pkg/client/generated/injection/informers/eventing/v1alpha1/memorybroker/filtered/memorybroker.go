// Copyright 2022 TriggerMesh Inc.
// SPDX-License-Identifier: Apache-2.0
// Code generated by injection-gen. DO NOT EDIT.

package filtered

import (
	context "context"

	apiseventingv1alpha1 "github.com/triggermesh/triggermesh-core/pkg/apis/eventing/v1alpha1"
	internalclientset "github.com/triggermesh/triggermesh-core/pkg/client/generated/clientset/internalclientset"
	v1alpha1 "github.com/triggermesh/triggermesh-core/pkg/client/generated/informers/externalversions/eventing/v1alpha1"
	client "github.com/triggermesh/triggermesh-core/pkg/client/generated/injection/client"
	filtered "github.com/triggermesh/triggermesh-core/pkg/client/generated/injection/informers/factory/filtered"
	eventingv1alpha1 "github.com/triggermesh/triggermesh-core/pkg/client/generated/listers/eventing/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
)

func init() {
	injection.Default.RegisterFilteredInformers(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct {
	Selector string
}

func withInformer(ctx context.Context) (context.Context, []controller.Informer) {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	infs := []controller.Informer{}
	for _, selector := range labelSelectors {
		f := filtered.Get(ctx, selector)
		inf := f.Eventing().V1alpha1().MemoryBrokers()
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
		infs = append(infs, inf.Informer())
	}
	return ctx, infs
}

func withDynamicInformer(ctx context.Context) context.Context {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	for _, selector := range labelSelectors {
		inf := &wrapper{client: client.Get(ctx), selector: selector}
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
	}
	return ctx
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context, selector string) v1alpha1.MemoryBrokerInformer {
	untyped := ctx.Value(Key{Selector: selector})
	if untyped == nil {
		logging.FromContext(ctx).Panicf(
			"Unable to fetch github.com/triggermesh/triggermesh-core/pkg/client/generated/informers/externalversions/eventing/v1alpha1.MemoryBrokerInformer with selector %s from context.", selector)
	}
	return untyped.(v1alpha1.MemoryBrokerInformer)
}

type wrapper struct {
	client internalclientset.Interface

	namespace string

	selector string
}

var _ v1alpha1.MemoryBrokerInformer = (*wrapper)(nil)
var _ eventingv1alpha1.MemoryBrokerLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apiseventingv1alpha1.MemoryBroker{}, 0, nil)
}

func (w *wrapper) Lister() eventingv1alpha1.MemoryBrokerLister {
	return w
}

func (w *wrapper) MemoryBrokers(namespace string) eventingv1alpha1.MemoryBrokerNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace, selector: w.selector}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apiseventingv1alpha1.MemoryBroker, err error) {
	reqs, err := labels.ParseToRequirements(w.selector)
	if err != nil {
		return nil, err
	}
	selector = selector.Add(reqs...)
	lo, err := w.client.EventingV1alpha1().MemoryBrokers(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apiseventingv1alpha1.MemoryBroker, error) {
	// TODO(mattmoor): Check that the fetched object matches the selector.
	return w.client.EventingV1alpha1().MemoryBrokers(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}
