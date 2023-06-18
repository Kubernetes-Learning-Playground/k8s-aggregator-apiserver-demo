/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/myoperator/k8saggregatorapiserver/pkg/apis/myingress/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMyIngresses implements MyIngressInterface
type FakeMyIngresses struct {
	Fake *FakeApisV1beta1
	ns   string
}

var myingressesResource = schema.GroupVersionResource{Group: "apis.jtthink.com", Version: "v1beta1", Resource: "myingresses"}

var myingressesKind = schema.GroupVersionKind{Group: "apis.jtthink.com", Version: "v1beta1", Kind: "MyIngress"}

// Get takes name of the myIngress, and returns the corresponding myIngress object, and an error if there is any.
func (c *FakeMyIngresses) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.MyIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(myingressesResource, c.ns, name), &v1beta1.MyIngress{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.MyIngress), err
}

// List takes label and field selectors, and returns the list of MyIngresses that match those selectors.
func (c *FakeMyIngresses) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.MyIngressList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(myingressesResource, myingressesKind, c.ns, opts), &v1beta1.MyIngressList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.MyIngressList{ListMeta: obj.(*v1beta1.MyIngressList).ListMeta}
	for _, item := range obj.(*v1beta1.MyIngressList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested myIngresses.
func (c *FakeMyIngresses) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(myingressesResource, c.ns, opts))

}

// Create takes the representation of a myIngress and creates it.  Returns the server's representation of the myIngress, and an error, if there is any.
func (c *FakeMyIngresses) Create(ctx context.Context, myIngress *v1beta1.MyIngress, opts v1.CreateOptions) (result *v1beta1.MyIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(myingressesResource, c.ns, myIngress), &v1beta1.MyIngress{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.MyIngress), err
}

// Update takes the representation of a myIngress and updates it. Returns the server's representation of the myIngress, and an error, if there is any.
func (c *FakeMyIngresses) Update(ctx context.Context, myIngress *v1beta1.MyIngress, opts v1.UpdateOptions) (result *v1beta1.MyIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(myingressesResource, c.ns, myIngress), &v1beta1.MyIngress{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.MyIngress), err
}

// Delete takes name of the myIngress and deletes it. Returns an error if one occurs.
func (c *FakeMyIngresses) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(myingressesResource, c.ns, name), &v1beta1.MyIngress{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMyIngresses) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(myingressesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.MyIngressList{})
	return err
}

// Patch applies the patch and returns the patched myIngress.
func (c *FakeMyIngresses) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.MyIngress, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(myingressesResource, c.ns, name, pt, data, subresources...), &v1beta1.MyIngress{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.MyIngress), err
}