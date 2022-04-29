/*
Copyright 2022.

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

	"strconv"

	"github.com/go-logr/logr"
	"github.com/prometheus/common/log"
	a "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "WeMainnet/api/v1"
)

// WeMainnetReconciler reconciles a WeMainnet object
type WeMainnetReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *WeMainnetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("WeMainnet", req.NamespacedName)
	log.Info("Processing WeMainnetReconciler.")
	WeMainnet := &appsv1.WeMainnet{}
	err := r.Client.Get(ctx, req.NamespacedName, WeMainnet)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("WeMainnet resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get WeMainnet")
		return ctrl.Result{}, err
	}
	// Check if the StatefulSet already exists, if not create a new one
	found := &a.StatefulSet{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: "node", Namespace: WeMainnet.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deployWeMainnetNode(WeMainnet)
		log.Info("Creating a new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", "node")
		err = r.Client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", "node")
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get StatefulSet")
		return ctrl.Result{}, err
	}

	// Check desired amount of sts.
	size := WeMainnet.Spec.ReplicasNode
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Client.Update(ctx, found)
		if err != nil {
			log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", found.Namespace, "StatefulSet.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}
	// Service for each NODE
	for i := 0; i < int(WeMainnet.Spec.ReplicasNode); i++ {
		foundServiceNode1 := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "node-" + strconv.Itoa(i), Namespace: found.Namespace}, foundServiceNode1)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceForEachNode(WeMainnet, int(i))
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "node-"+strconv.Itoa(i))
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "node-"+strconv.Itoa(i))
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}

	// NODE BALANCER
	foundServiceNode := &corev1.Service{}
	// repl := ha.Spec.Replicas
	// print()
	err = r.Client.Get(ctx, types.NamespacedName{Name: "node", Namespace: found.Namespace}, foundServiceNode)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		dep := r.serviceForNode(WeMainnet)
		log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "node")
		err = r.Client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "node")
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	//Dockerhost
	foundDockerhost := &a.Deployment{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: "dockerhost", Namespace: WeMainnet.Namespace}, foundDockerhost)
	if err != nil && errors.IsNotFound(err) {
		dep := r.deployWeMainnetDockerhost(WeMainnet)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "dockerhost")
		err = r.Client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "dockerhost")
			return ctrl.Result{}, err
		}
		// StatefulSet created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	sizeDockerhost := int32(1)
	if *foundDockerhost.Spec.Replicas != sizeDockerhost {
		foundDockerhost.Spec.Replicas = &sizeDockerhost
		err = r.Client.Update(ctx, foundDockerhost)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDockerhost.Namespace, "Deployment.Name", foundDockerhost.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}
	//SERVICE DOCKERHOST
	foundServiceDockerhost := &corev1.Service{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: "dockerhost", Namespace: found.Namespace}, foundServiceDockerhost)
	if err != nil && errors.IsNotFound(err) {
		// Define a new service
		dep := r.serviceForDockerhost(WeMainnet)
		log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "dockerhost")
		err = r.Client.Create(ctx, dep)
		if err != nil {
			log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceDockerhost")
			return ctrl.Result{}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{}, err
	}

	if WeMainnet.Spec.TelegrafEnable == "" {
		WeMainnet.Spec.TelegrafEnable = "false"
	}
	if WeMainnet.Spec.TelegrafEnable == "true" {
		//Telegraf
		foundTelegraf := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "telegraf", Namespace: WeMainnet.Namespace}, foundTelegraf)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetTelegraf(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "telegraf")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "telegraf")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeTelegraf := WeMainnet.Spec.ReplicasTelegraf
		if *foundTelegraf.Spec.Replicas != sizeTelegraf {
			foundTelegraf.Spec.Replicas = &sizeTelegraf
			err = r.Client.Update(ctx, foundTelegraf)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundTelegraf.Namespace, "Deployment.Name", foundTelegraf.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}
		// Service telegraf
		foundServiceTelegraf := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "telegraf", Namespace: found.Namespace}, foundServiceTelegraf)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceTelegraf(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "telegraf")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceTelegraf")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}

	if WeMainnet.Spec.AuthAdminEnable == "" {
		WeMainnet.Spec.AuthAdminEnable = "false"
	}
	if WeMainnet.Spec.AuthAdminEnable == "true" {
		//Auth-Admin
		foundAuthAdmin := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "auth-admin", Namespace: WeMainnet.Namespace}, foundAuthAdmin)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetAuthAdmin(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "auth-admin")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "auth-admin")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeAuthAdmin := WeMainnet.Spec.ReplicasAuthAdmin
		if *foundAuthAdmin.Spec.Replicas != sizeAuthAdmin {
			foundAuthAdmin.Spec.Replicas = &sizeAuthAdmin
			err = r.Client.Update(ctx, foundAuthAdmin)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundAuthAdmin.Namespace, "Deployment.Name", foundAuthAdmin.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}
		// Service auth-admin
		foundServiceAuthAdmin := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "auth-admin", Namespace: found.Namespace}, foundServiceAuthAdmin)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceAuthAdmin(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "auth-admin")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceAuthAdmin")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}
	if WeMainnet.Spec.PostgresEnable == "true" {
		//Postgresql
		foundPostgresql := &a.StatefulSet{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "postgresql", Namespace: WeMainnet.Namespace}, foundPostgresql)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetPostgresql(WeMainnet)
			log.Info("Creating a new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", "postgresql")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new StatefulSet", "StatefulSet.Namespace", dep.Namespace, "StatefulSet.Name", "postgresql")
				return ctrl.Result{}, err
			}
			// StatefulSet created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get StatefulSet")
			return ctrl.Result{}, err
		}

		sizePostgresql := int32(1)
		if *foundPostgresql.Spec.Replicas != sizePostgresql {
			foundPostgresql.Spec.Replicas = &sizePostgresql
			err = r.Client.Update(ctx, foundPostgresql)
			if err != nil {
				log.Error(err, "Failed to update StatefulSet", "StatefulSet.Namespace", foundPostgresql.Namespace, "StatefulSet.Name", foundPostgresql.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		//Service postgresql
		foundServicePostgresql := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "postgresql", Namespace: found.Namespace}, foundServicePostgresql)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceForPostgresql(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "postgresql")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "postgresql")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}
	if WeMainnet.Spec.ClientEnable == "true" {
		//Auth
		foundAuth := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "auth-service", Namespace: WeMainnet.Namespace}, foundAuth)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetAuth(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "auth-service")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "auth-service")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeAuth := WeMainnet.Spec.ReplicasAuth
		if *foundAuth.Spec.Replicas != sizeAuth {
			foundAuth.Spec.Replicas = &sizeAuth
			err = r.Client.Update(ctx, foundAuth)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundAuth.Namespace, "Deployment.Name", foundAuth.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		//Nginx
		foundNginx := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "nginx", Namespace: WeMainnet.Namespace}, foundNginx)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployNginx(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "nginx")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "nginx")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeNginx := WeMainnet.Spec.ReplicasNginx
		if *foundNginx.Spec.Replicas != sizeNginx {
			foundNginx.Spec.Replicas = &sizeNginx
			err = r.Client.Update(ctx, foundNginx)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundNginx.Namespace, "Deployment.Name", foundNginx.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		// Service auth
		foundServiceAuth := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "service-auth-service", Namespace: found.Namespace}, foundServiceAuth)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceForAuth(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "service-auth-service")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "service-auth-service")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}

		//Crawler
		foundCrawler := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "crawler2", Namespace: WeMainnet.Namespace}, foundCrawler)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetCrawler(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "crawler2")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "crawler2")
				return ctrl.Result{}, err
			}
			// StatefulSet created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get StatefulSet")
			return ctrl.Result{}, err
		}

		sizeCrawler := WeMainnet.Spec.ReplicasCrawler
		if *foundCrawler.Spec.Replicas != sizeCrawler {
			foundCrawler.Spec.Replicas = &sizeCrawler
			err = r.Client.Update(ctx, foundCrawler)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundCrawler.Namespace, "Deployment.Name", foundCrawler.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		//Dataservice
		foundDs := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "dataservice", Namespace: WeMainnet.Namespace}, foundDs)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetDs(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "dataservice")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "dataservice")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeDs := WeMainnet.Spec.ReplicasDs
		if *foundDs.Spec.Replicas != sizeDs {
			foundDs.Spec.Replicas = &sizeDs
			err = r.Client.Update(ctx, foundDs)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDs.Namespace, "Deployment.Name", foundDs.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		//Frontend
		foundFrontend := &a.Deployment{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "frontend", Namespace: WeMainnet.Namespace}, foundFrontend)
		if err != nil && errors.IsNotFound(err) {
			dep := r.deployWeMainnetFrontend(WeMainnet)
			log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "frontend")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", "frontend")
				return ctrl.Result{}, err
			}
			// Deployment created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Deployment")
			return ctrl.Result{}, err
		}

		sizeFrontend := WeMainnet.Spec.ReplicasFrontend
		if *foundFrontend.Spec.Replicas != sizeFrontend {
			foundFrontend.Spec.Replicas = &sizeFrontend
			err = r.Client.Update(ctx, foundFrontend)
			if err != nil {
				log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDs.Namespace, "Deployment.Name", foundFrontend.Name)
				return ctrl.Result{}, err
			}
			// Spec updated - return and requeue
			return ctrl.Result{Requeue: true}, nil
		}

		// Service frontend
		foundServiceFrontend := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "frontend", Namespace: found.Namespace}, foundServiceFrontend)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceForFrontend(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "frontend")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceFrontend")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}

		// Service ds
		foundServiceDs := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "dataservice", Namespace: found.Namespace}, foundServiceDs)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceForDs(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "dataservice")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceDataservice")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}

		// Service nginx
		foundServiceNginx := &corev1.Service{}
		err = r.Client.Get(ctx, types.NamespacedName{Name: "nginx", Namespace: found.Namespace}, foundServiceNginx)
		if err != nil && errors.IsNotFound(err) {
			// Define a new service
			dep := r.serviceNginx(WeMainnet)
			log.Info("Creating a new Service", "Service.Namespace", dep.Namespace, "Service.Name", "nginx")
			err = r.Client.Create(ctx, dep)
			if err != nil {
				log.Error(err, "Failed to create new Service", "Service.Namespace", dep.Namespace, "Service.Name", "foundServiceNginx")
				return ctrl.Result{}, err
			}
			// Service created successfully - return and requeue
			return ctrl.Result{Requeue: true}, nil
		} else if err != nil {
			log.Error(err, "Failed to get Service")
			return ctrl.Result{}, err
		}
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WeMainnetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.WeMainnet{}).
		Complete(r)
}
