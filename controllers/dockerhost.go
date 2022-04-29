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

//Dockerhost
func (c *WeMainnetReconciler) deployWeMainnetDockerhost(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := int32(1)
	labels := map[string]string{"app": "dockerhost"}
	image := ha.Spec.ImageDockerhost
	//not defined vars
	if ha.Spec.CpuDockerhostRequest == "" {
		ha.Spec.CpuDockerhostRequest = "1"
	}
	if ha.Spec.CpuDockerhostLimit == "" {
		ha.Spec.CpuDockerhostLimit = "2"
	}
	if ha.Spec.MemoryDockerhostRequest == "" {
		ha.Spec.MemoryDockerhostRequest = "1Gi"
	}
	if ha.Spec.MemoryDockerhostLimit == "" {
		ha.Spec.MemoryDockerhostLimit = "2Gi"
	}

	//
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dockerhost",
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
							Name: "varlibdocker",
							VolumeSource: corev1.VolumeSource{
								EmptyDir: &corev1.EmptyDirVolumeSource{},
							}},
					},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "dockerhost",
						SecurityContext: &corev1.SecurityContext{
							Privileged: &ha.Spec.PrivilegedDockerhost,
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryDockerhostRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuDockerhostRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryDockerhostLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuDockerhostLimit),
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "DOCKER_HOST",
								Value: "tcp://localhost:2375",
							},
							{
								Name:  "DOCKER_TLS_CERTDIR",
								Value: "",
							}},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "varlibdocker",
								MountPath: "/var/lib/docker",
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

//Dockerhost service
func (r *WeMainnetReconciler) serviceForDockerhost(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "dockerhost"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dockerhost",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 2375,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 2375,
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
