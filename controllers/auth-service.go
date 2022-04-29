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

//AUTH
func (c *WeMainnetReconciler) deployWeMainnetAuth(ha *appsv1.WeMainnet) *a.Deployment {
	replicas := ha.Spec.ReplicasAuth
	labels := map[string]string{"app": "auth-service"}
	image := ha.Spec.ImageAuth

	if ha.Spec.MailEnabled == "" {
		ha.Spec.MailEnabled = "false"
	}
	if ha.Spec.ActivateUserOnRegister == "" {
		ha.Spec.ActivateUserOnRegister = "true"
	}
	if ha.Spec.CpuAuthRequest == "" {
		ha.Spec.CpuAuthRequest = "100m"
	}
	if ha.Spec.CpuAuthLimit == "" {
		ha.Spec.CpuAuthLimit = "1"
	}
	if ha.Spec.MemoryAuthRequest == "" {
		ha.Spec.MemoryAuthRequest = "32Mi"
	}
	if ha.Spec.MemoryAuthLimit == "" {
		ha.Spec.MemoryAuthLimit = "1Gi"
	}
	//
	dep := &a.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "auth-service",
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
							Name: "auth-service-keys",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "auth-service-keys",
								},
							},
						},
						{
							Name: "auth-service-tokens",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "auth-service-tokens",
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "auth-service",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryAuthRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuAuthRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryAuthLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuAuthLimit),
							},
						},
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: 10,
							PeriodSeconds:       5,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "status",
									Port:   intstr.FromInt(3000),
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
									Path:   "status",
									Port:   intstr.FromInt(3000),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "auth-service-keys",
								ReadOnly:  true,
								MountPath: "/etc/auth-service-keys",
							},
							{
								Name:      "auth-service-tokens",
								ReadOnly:  true,
								MountPath: "/app/tokens.json",
								SubPath:   "tokens.json",
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "JWT_EXPIRATION_DATE",
								Value: "3600s",
							},
							{
								Name:  "REFRESH_EXPIRATION_DATE",
								Value: "3600s",
							},
							{
								Name:  "ACTIVATE_USER_ON_REGISTER",
								Value: ha.Spec.ActivateUserOnRegister,
							},
							{
								Name:  "IS_MAIL_TRANSPORT_ENABLED",
								Value: ha.Spec.MailEnabled,
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
								Value: "auth_service",
							},
							{
								Name: "MAIL_HOST",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-host",
									},
								},
							},
							{
								Name: "MAIL_USER",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-user",
									},
								},
							},
							{
								Name: "MAIL_PASSWORD",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-password",
									},
								},
							},
							{
								Name: "MAIL_PORT",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-port",
									},
								},
							},
							{
								Name: "MAIL_FROM",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-from",
									},
								},
							},
							{
								Name: "PASSWORD_HASH_SALT",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "mail",
										},
										Key: "mail-salt",
									},
								},
							},
							{
								Name:  "RSA_PRIVATE_FILE_PATH",
								Value: "/etc/auth-service-keys/jwtRS256.key",
							},
							{
								Name:  "RSA_PUBLIC_FILE_PATH",
								Value: "/etc/auth-service-keys/jwtRS256.key.pub",
							}},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(ha, dep, c.Scheme)
	return dep
}

//Auth Service
func (r *WeMainnetReconciler) serviceForAuth(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "auth-service"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "service-auth-service",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 3000,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 3000,
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
