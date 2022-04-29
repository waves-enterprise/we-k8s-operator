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

//Frontend
func (c *WeMainnetReconciler) deployWeMainnetFrontend(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasFrontend
	labels := map[string]string{"app": "frontend"}
	image := ha.Spec.ImageFrontend
	//not defined vars
	if ha.Spec.CpuFrontendRequest == "" {
		ha.Spec.CpuFrontendRequest = "1"
	}
	if ha.Spec.CpuFrontendLimit == "" {
		ha.Spec.CpuFrontendLimit = "1"
	}
	if ha.Spec.MemoryFrontendRequest == "" {
		ha.Spec.MemoryFrontendRequest = "1Gi"
	}
	if ha.Spec.MemoryFrontendLimit == "" {
		ha.Spec.MemoryFrontendLimit = "1Gi"
	}
	//
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend",
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
					Volumes: []corev1.Volume{
						{
							Name: "frontend-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "frontend-config",
									},
								},
							},
						}},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "frontend",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryFrontendRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuFrontendRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryFrontendLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuFrontendLimit),
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(8080),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(8080),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "frontend-config",
								SubPath:   "app.config.json",
								MountPath: "/usr/share/nginx/html/app.config.json",
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

//Frontend service
func (r *WeMainnetReconciler) serviceForFrontend(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "frontend"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "frontend",
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
