/*
Copyright 2019 BlackRock, Inc.

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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/beval/beval/pkg/apis/beval/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// bevalLister helps list bevals.
type bevalLister interface {
	// List lists all bevals in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.beval, err error)
	// bevals returns an object that can list and get bevals.
	bevals(namespace string) bevalNamespaceLister
	bevalListerExpansion
}

// bevalLister implements the bevalLister interface.
type bevalLister struct {
	indexer cache.Indexer
}

// NewbevalLister returns a new bevalLister.
func NewbevalLister(indexer cache.Indexer) bevalLister {
	return &bevalLister{indexer: indexer}
}

// List lists all bevals in the indexer.
func (s *bevalLister) List(selector labels.Selector) (ret []*v1alpha1.beval, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.beval))
	})
	return ret, err
}

// bevals returns an object that can list and get bevals.
func (s *bevalLister) bevals(namespace string) bevalNamespaceLister {
	return bevalNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// bevalNamespaceLister helps list and get bevals.
type bevalNamespaceLister interface {
	// List lists all bevals in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.beval, err error)
	// Get retrieves the beval from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.beval, error)
	bevalNamespaceListerExpansion
}

// bevalNamespaceLister implements the bevalNamespaceLister
// interface.
type bevalNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all bevals in the indexer for a given namespace.
func (s bevalNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.beval, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.beval))
	})
	return ret, err
}

// Get retrieves the beval from the indexer for a given namespace and name.
func (s bevalNamespaceLister) Get(name string) (*v1alpha1.beval, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("beval"), name)
	}
	return obj.(*v1alpha1.beval), nil
}
