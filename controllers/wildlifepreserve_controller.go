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

package controllers

import (
	"context"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	wildlifeV1 "github.com/jakubpieta/wildlife-preserve-operator/api/v1alpha1"
)

// WildlifePreserveReconciler reconciles a WildlifePreserve object
type WildlifePreserveReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=wildlife.preserves.jakubpieta,resources=wildlifepreserves,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumes,verbs=get;list;create;update;patch;delete
//+kubebuilder:rbac:groups=wildlife.preserves.jakubpieta,resources=wildlifepreserves/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=wildlife.preserves.jakubpieta,resources=wildlifepreserves/finalizers,verbs=update

// Reconcile handles WildlifePreserve resource reconciliation
func (r *WildlifePreserveReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Fetch the WildlifePreserve resource
	preserve := &wildlifeV1.WildlifePreserve{}
	if err := r.Get(ctx, req.NamespacedName, preserve); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Handle the volume creation based on the `Volume` field
	if preserve.Spec.VolumeMountPath != "" {
		pvc := &coreV1.PersistentVolumeClaim{}
		err := r.Get(ctx, types.NamespacedName{
			Namespace: req.Namespace,
			Name:      preserve.Name + "-volume-claim",
		}, pvc)
		if errors.IsNotFound(err) {
			// Create a volume with the specified configuration
			// You can use Kubernetes client (client.Client) to create the volume (e.g., PVC)
			// Example: Create a PVC (PersistentVolumeClaim)
			pvc = &coreV1.PersistentVolumeClaim{
				ObjectMeta: metaV1.ObjectMeta{
					Name:      preserve.Name + "-volume-claim",
					Namespace: req.Namespace,
				},
				Spec: coreV1.PersistentVolumeClaimSpec{
					AccessModes: []coreV1.PersistentVolumeAccessMode{coreV1.ReadWriteOnce},
					Resources: coreV1.ResourceRequirements{
						Requests: coreV1.ResourceList{
							coreV1.ResourceStorage: resource.MustParse("100Mi"), // Adjust the storage size as needed
						},
					},
				},
			}
			err = r.Create(ctx, pvc)
			if err != nil && !errors.IsAlreadyExists(err) {
				return ctrl.Result{}, err
			}
		}
	}

	// Create or update a Deployment for the Go application with volume mount
	deployment := &appsV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name:      preserve.Name + "-deployment",
			Namespace: req.Namespace,
		},
		Spec: appsV1.DeploymentSpec{
			Replicas: &preserve.Spec.Replicas, // Use the number of replicas from the CR
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{"app": preserve.Name},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Labels: map[string]string{"app": preserve.Name},
				},
				Spec: coreV1.PodSpec{
					Volumes: []coreV1.Volume{
						{
							Name: preserve.Name + "-volume",
							VolumeSource: coreV1.VolumeSource{
								PersistentVolumeClaim: &coreV1.PersistentVolumeClaimVolumeSource{
									ClaimName: preserve.Name + "-volume-claim",
								},
							},
						},
					},
					Containers: []coreV1.Container{
						{
							Name:  "wildlife-preserve-app",
							Image: "jakubpieta/wildlife-preserve-app:v0.0.1", // Replace with your Go app image
							VolumeMounts: []coreV1.VolumeMount{
								{
									Name:      preserve.Name + "-volume",
									MountPath: preserve.Spec.VolumeMountPath,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create or update the Deployment
	if err := controllerutil.SetControllerReference(preserve, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	err := r.Create(ctx, deployment)
	if err != nil && !errors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	// Handle other aspects of the WildlifePreserve resource

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WildlifePreserveReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&wildlifeV1.WildlifePreserve{}).
		Complete(r)
}
