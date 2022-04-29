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

//Telegraf
func (c *WeMainnetReconciler) deployWeMainnetTelegraf(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasTelegraf
	labels := map[string]string{"app": "telegraf"}
	image := ha.Spec.ImageTelegraf
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "telegraf",
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
					Volumes: []corev1.Volume{{
						Name: "telegraf-config",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "telegraf-config",
								},
							},
						},
					}},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "telegraf",
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
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "telegraf-config",
								MountPath: "/etc/telegraf",
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

//Telegraf service
func (r *WeMainnetReconciler) serviceTelegraf(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "telegraf"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "telegraf",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 8086,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 8086,
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
