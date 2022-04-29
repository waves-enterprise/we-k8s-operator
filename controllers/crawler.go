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

//CRAWLER
func (c *WeMainnetReconciler) deployWeMainnetCrawler(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasCrawler
	labels := map[string]string{"app": "crawler2"}
	image := ha.Spec.ImageCrawler

	//not defined vars
	if ha.Spec.CrawlerServiceToken == "" {
		ha.Spec.CrawlerServiceToken = "Mh7K0iMj1je3prC3kkjhZE3FOX5inPRbXvOAhsPR"
	}
	if ha.Spec.GrpcAddresses == "" {
		ha.Spec.GrpcAddresses = "node-0:6865,node-1:6865,node-2:6865"
	}
	if ha.Spec.CpuCrawlerRequest == "" {
		ha.Spec.CpuCrawlerRequest = "1"
	}
	if ha.Spec.CpuCrawlerLimit == "" {
		ha.Spec.CpuCrawlerLimit = "1"
	}
	if ha.Spec.MemoryCrawlerRequest == "" {
		ha.Spec.MemoryCrawlerRequest = "2Gi"
	}
	if ha.Spec.MemoryCrawlerLimit == "" {
		ha.Spec.MemoryCrawlerLimit = "2Gi"
	}

	//
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "crawler2",
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
						Name:  "crawler2",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryCrawlerRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuCrawlerRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryCrawlerLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuCrawlerLimit),
							},
						},
						LivenessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "livenessProbe",
									Port:   intstr.FromInt(3000),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "readinessProbe",
									Port:   intstr.FromInt(3000),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "AUTH_SERVICE_ADDRESS",
								Value: "http://service-auth-service:3000",
							},
							{
								Name:  "NODE_ADDRESS",
								Value: "http://node:6862",
							},
							{
								Name:  "SERVICE_TOKEN",
								Value: ha.Spec.CrawlerServiceToken,
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
							{
								Name:  "ROLLBACK_COUNT",
								Value: "100",
							},
							{
								Name:  "MAX_ASYNC_BLOCKS_REQUESTS",
								Value: "2",
							},
							{
								Name:  "MAX_BLOCKS_SIZE",
								Value: "5",
							},
							{
								Name:  "GRPC_ADDRESSES",
								Value: ha.Spec.GrpcAddresses,
							},
							{
								Name:  "NODE_API_KEY",
								Value: "",
							},
							{
								Name:  "ASYNC_GRPC",
								Value: "false",
							},
							{
								Name:  "NODE_JS_RAM_LIMIT",
								Value: "1024",
							}},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(ha, dep, c.Scheme)
	return dep
}
