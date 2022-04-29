package controllers

import (
	"strconv"

	a "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"

	appsv1 "WeMainnet/api/v1"
)

func (c *WeMainnetReconciler) deployWeMainnetNode(ha *appsv1.WeMainnet) *a.StatefulSet {
	replicas := ha.Spec.ReplicasNode
	labels := map[string]string{"app": "node"}
	image := ha.Spec.ImageNode
	//not defined vars
	if ha.Spec.JavaOpts == "" {
		ha.Spec.JavaOpts = "-Dwe.check-resources=false -Xmx3g"
	}
	if ha.Spec.CleanState == "" {
		ha.Spec.CleanState = "false"
	}
	if ha.Spec.CpuNodeRequest == "" {
		ha.Spec.CpuNodeRequest = "2"
	}
	if ha.Spec.CpuNodeLimit == "" {
		ha.Spec.CpuNodeLimit = "2"
	}
	if ha.Spec.MemoryNodeRequest == "" {
		ha.Spec.MemoryNodeRequest = "4Gi"
	}
	if ha.Spec.MemoryNodeLimit == "" {
		ha.Spec.MemoryNodeLimit = "4Gi"
	}
	//
	dep := &a.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
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
						Name: "data-volume",
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
							Name: "node-config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "node-config",
									},
								},
							},
						},
						{
							Name: "node-wallet",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: "node-wallet",
								},
							},
						},
						{
							Name: "node-license",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "node-license",
									},
								},
							},
						},
					},
					Containers: []corev1.Container{{
						Image: image,
						Name:  "node",
						Resources: corev1.ResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryNodeRequest),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuNodeRequest),
							},
							Limits: corev1.ResourceList{
								corev1.ResourceMemory: resource.MustParse(ha.Spec.MemoryNodeLimit),
								corev1.ResourceCPU:    resource.MustParse(ha.Spec.CpuNodeLimit),
							},
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "data-volume",
								MountPath: "/node/data",
							},
							{
								Name:      "node-config",
								MountPath: "/opt/configs",
							},
							{
								Name:      "node-wallet",
								MountPath: "/opt/wallets",
							},
							{
								Name:      "node-license",
								MountPath: "/opt/licenses",
							},
						},
						LivenessProbe: &corev1.Probe{
							InitialDelaySeconds: 30,
							PeriodSeconds:       20,
							FailureThreshold:    3,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "node/status",
									Port:   intstr.FromInt(6862),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						ReadinessProbe: &corev1.Probe{
							InitialDelaySeconds: 30,
							PeriodSeconds:       20,
							TimeoutSeconds:      10,
							Handler: corev1.Handler{
								HTTPGet: &corev1.HTTPGetAction{
									Path:   "node/status",
									Port:   intstr.FromInt(6862),
									Scheme: corev1.URISchemeHTTP,
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "vars-node",
									},
								},
							},
						},
						Env: []corev1.EnvVar{
							{
								Name:  "JAVA_OPTS",
								Value: ha.Spec.JavaOpts,
							},
							{
								Name:  "CONFIGNAME_AS_HOSTNAME",
								Value: "true",
							},
							{
								Name:  "CLEAN_STATE",
								Value: ha.Spec.CleanState,
							}},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(ha, dep, c.Scheme)
	return dep
}

func (r *WeMainnetReconciler) serviceForEachNode(ha *appsv1.WeMainnet, i int) *corev1.Service {
	ls := map[string]string{"app": "node"}
	sel := map[string]string{"statefulset.kubernetes.io/pod-name": "node-" + strconv.Itoa(i)}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node-" + strconv.Itoa(i),
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: sel,
			Ports: []corev1.ServicePort{
				{
					Port: 6862,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6862,
						Type:   intstr.Int,
					},
				},
				{
					Port: 6864,
					Name: "data-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6864,
						Type:   intstr.Int,
					},
				},
				{
					Port: 6865,
					Name: "grpc-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6865,
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

func (r *WeMainnetReconciler) serviceForNode(ha *appsv1.WeMainnet) *corev1.Service {

	ls := map[string]string{"app": "node"}

	dep := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "node",
			Namespace: ha.Namespace,
			Labels:    ls,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{
				{
					Port: 6862,
					Name: "http-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6862,
						Type:   intstr.Int,
					},
				},
				{
					Port: 6864,
					Name: "data-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6864,
						Type:   intstr.Int,
					},
				},
				{
					Port: 6865,
					Name: "grpc-port",
					TargetPort: intstr.IntOrString{
						IntVal: 6865,
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
