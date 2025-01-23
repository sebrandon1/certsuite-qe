// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	v1 "k8s.io/client-go/applyconfigurations/meta/v1"
)

// ObjectMetaApplyConfiguration represents a declarative configuration of the ObjectMeta type for use
// with apply.
type ObjectMetaApplyConfiguration struct {
	Name            *string                               `json:"name,omitempty"`
	GenerateName    *string                               `json:"generateName,omitempty"`
	Namespace       *string                               `json:"namespace,omitempty"`
	Labels          map[string]string                     `json:"labels,omitempty"`
	Annotations     map[string]string                     `json:"annotations,omitempty"`
	OwnerReferences []v1.OwnerReferenceApplyConfiguration `json:"ownerReferences,omitempty"`
}

// ObjectMetaApplyConfiguration constructs a declarative configuration of the ObjectMeta type for use with
// apply.
func ObjectMeta() *ObjectMetaApplyConfiguration {
	return &ObjectMetaApplyConfiguration{}
}

// WithName sets the Name field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Name field is set to the value of the last call.
func (b *ObjectMetaApplyConfiguration) WithName(value string) *ObjectMetaApplyConfiguration {
	b.Name = &value
	return b
}

// WithGenerateName sets the GenerateName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the GenerateName field is set to the value of the last call.
func (b *ObjectMetaApplyConfiguration) WithGenerateName(value string) *ObjectMetaApplyConfiguration {
	b.GenerateName = &value
	return b
}

// WithNamespace sets the Namespace field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Namespace field is set to the value of the last call.
func (b *ObjectMetaApplyConfiguration) WithNamespace(value string) *ObjectMetaApplyConfiguration {
	b.Namespace = &value
	return b
}

// WithLabels puts the entries into the Labels field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Labels field,
// overwriting an existing map entries in Labels field with the same key.
func (b *ObjectMetaApplyConfiguration) WithLabels(entries map[string]string) *ObjectMetaApplyConfiguration {
	if b.Labels == nil && len(entries) > 0 {
		b.Labels = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Labels[k] = v
	}
	return b
}

// WithAnnotations puts the entries into the Annotations field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, the entries provided by each call will be put on the Annotations field,
// overwriting an existing map entries in Annotations field with the same key.
func (b *ObjectMetaApplyConfiguration) WithAnnotations(entries map[string]string) *ObjectMetaApplyConfiguration {
	if b.Annotations == nil && len(entries) > 0 {
		b.Annotations = make(map[string]string, len(entries))
	}
	for k, v := range entries {
		b.Annotations[k] = v
	}
	return b
}

// WithOwnerReferences adds the given value to the OwnerReferences field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the OwnerReferences field.
func (b *ObjectMetaApplyConfiguration) WithOwnerReferences(values ...*v1.OwnerReferenceApplyConfiguration) *ObjectMetaApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithOwnerReferences")
		}
		b.OwnerReferences = append(b.OwnerReferences, *values[i])
	}
	return b
}
