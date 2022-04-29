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

//Postgresql
func (c *WeMainnetReconciler) deployWeMainnetPostgresql(ha *appsv1.WeMainnet) *a.StatefulSet {
	replicas := int32(1)
	labels := map[string]string{"app": "postgresql"}
	image := "postgres:10-alpine"

	if ha.Spec.CpuAuthRequest == "" {
		ha.Spec.CpuAuthRequest = "500m"
	}
	if ha.Spec.CpuAuthLimit == "" {
		ha.Spec.CpuAuthLimit = "2"
	}
	if ha.Spec.MemoryAuthRequest == "" {
		ha.Spec.MemoryAuthRequest = "512Mi"
	}
	if ha.Spec.MemoryAuthLimit == "" {
		ha.Spec.MemoryAuthLimit = "2Gi"
	}
	//
	dep := &a.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: ha.Namespace,
		},
		Spec: a.StatefulSetSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "dbdata",
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							"ReadWriteOnce",
						},
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse(ha.Spec.Storage),
							},
						},
					},
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Volumes: []corev1.Volume{
						{
							Name: "postgresql",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "postgresql",
								},
							},
						},
						{
							Name: "init-script",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "init-script",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "postgresql",
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
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "dbdata",
								MountPath: "/var/lib/postgresql/data",
							},
							{
								Name:      "init-script",
								SubPath:   "create-multiple-postgresql-databases.sh",
								MountPath: "/docker-entrypoint-initdb.d/create-multiple-postgresql-databases.sh",
							},
						},
						Env: []corev1.EnvVar{
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
								Name:  "PGDATA",
								Value: "/var/lib/postgresql/data/pgdata",
							},
							{
								Name:  "POSTGRES_MULTIPLE_DATABASES",
								Value: "blockchain,auth_service",
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
func (r *WeMainnetReconciler) serviceForPostgresql(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "postgresql"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "postgresql",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 5432,
					Name: "tcp-port",
					TargetPort: intstr.IntOrString{
						IntVal: 5432,
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
