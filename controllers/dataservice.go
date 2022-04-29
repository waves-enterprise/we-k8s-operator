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

//DATASERVICE
func (c *WeMainnetReconciler) deployWeMainnetDs(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasDs
	labels := map[string]string{"app": "dataservice"}
	image := ha.Spec.ImageDs

	//not defined vars
	if ha.Spec.DsServiceToken == "" {
		ha.Spec.DsServiceToken = "Mh7K0iMj1je3prC3kkjhZE3FOX5inPRbXvOAhsPR"
	}
	if ha.Spec.CpuDsRequest == "" {
		ha.Spec.CpuDsRequest = "1"
	}
	if ha.Spec.CpuDsLimit == "" {
		ha.Spec.CpuDsLimit = "1"
	}
	if ha.Spec.MemoryDsRequest == "" {
		ha.Spec.MemoryDsRequest = "1Gi"
	}
	if ha.Spec.MemoryDsLimit == "" {
		ha.Spec.MemoryDsLimit = "1Gi"
	}
	//
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dataservice",
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
						Name:  "dataservice",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryDsRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuDsRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryDsLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuDsLimit),
							},
						},
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(3001),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
							SuccessThreshold:    1,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Port:   intstr.FromInt(3001),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "API_PORT",
								Value: "3001",
							},
							{
								Name:  "NODE_ENV",
								Value: "production",
							},
							{
								Name:  "LOG_LEVEL",
								Value: "info",
							},
							{
								Name:  "VOSTOK_NODE_ADDRESS",
								Value: "http://node:6862",
							},
							{
								Name:  "VOSTOK_AUTH_SERVICE_ADDRESS",
								Value: "http://service-auth-service:3000",
							},
							{
								Name:  "VOSTOK_AUTH_SERVICE_ADDRESS",
								Value: "http://service-auth-service:3000",
							},
							{
								Name:  "VOSTOK_AUTH_SERVICE_TOKEN",
								Value: ha.Spec.DsServiceToken,
							},
							{
								Name: "POSTGRES_ENABLE_SSL",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "pg-user",
										},
										Key: "ssl",
									},
								},
							},
							{
								Name:  "POSTGRES_HOST",
								Value: "postgresql",
							},
							{
								Name: "POSTGRES_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "pg-user",
										},
										Key: "user",
									},
								},
							},
							{
								Name: "POSTGRES_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "pg-user",
										},
										Key: "password_admin",
									},
								},
							},
							{
								Name:  "POSTGRES_DB",
								Value: "blockchain",
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

//Dataservice service
func (r *WeMainnetReconciler) serviceForDs(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "dataservice"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "dataservice",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 3001,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 3001,
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
