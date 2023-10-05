/*
Copyright 2023.

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

package v1alpha1

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WildlifePreserveSpec defines the desired state of WildlifePreserve
type WildlifePreserveSpec struct {
	Name            string `json:"name"`
	Location        string `json:"location"`
	Replicas        int32  `json:"replicas"`
	VolumeMountPath string `json:"volumeMountPath"`
}

// WildlifePreserveStatus defines the observed state of WildlifePreserve
type WildlifePreserveStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WildlifePreserve is the Schema for the wildlife-preserves API
type WildlifePreserve struct {
	metaV1.TypeMeta   `json:",inline"`
	metaV1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WildlifePreserveSpec   `json:"spec,omitempty"`
	Status WildlifePreserveStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WildlifePreserveList contains a list of WildlifePreserve
type WildlifePreserveList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata,omitempty"`
	Items           []WildlifePreserve `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WildlifePreserve{}, &WildlifePreserveList{})
}
