package controllers

import (
	a "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	appsv1 "WeMainnet/api/v1"
)

//Auth-admin
func (c *WeMainnetReconciler) deployWeMainnetAuthAdmin(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasAuthAdmin
	labels := map[string]string{"app": "auth-admin"}
	image := ha.Spec.ImageAuthAdmin
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auth-admin",
			Namespace: ha.Namespace,
		},
		Spec: a.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},

			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "auth-admin",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("100Mi"),
								corev1.ResourceCPU:    resource.MustParse("100m"),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse("256Mi"),
								corev1.ResourceCPU:    resource.MustParse("500m"),
							},
						},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(ha, dep, c.Scheme)
	return dep
}

//Auth-admin service
func (r *WeMainnetReconciler) serviceAuthAdmin(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "auth-admin"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auth-admin",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 8080,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 8080,
						Type:   intstr.Int,
					},
				},
			},
			Type: "ClusterIP",
		},
	}
	// Set Operator instance as the owner and controller
	ctrl.SetControllerReference(ha, dep, r.Scheme)
	return dep
}
