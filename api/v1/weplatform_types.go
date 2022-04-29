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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// WeMainnetSpec defines the desired state of WeMainnet
type WeMainnetSpec struct {
	ClientEnable    string `json:"client_enabled"`
	PostgresEnable  string `json:"postgres_enabled"`
	TelegrafEnable  string `json:"telegraf_enabled,omitempty"`
	AuthAdminEnable string `json:"authadmin_enabled,omitempty"`
	//NODE
	ImageNode         string `json:"image"`
	CleanState        string `json:"clean_state,omitempty"`
	ReplicasNode      int32  `json:"replicas"`
	CpuNodeRequest    string `json:"cpu_node_request,omitempty"`
	CpuNodeLimit      string `json:"cpu_node_limit,omitempty"`
	Storage           string `json:"storage"`
	MemoryNodeRequest string `json:"memory_node_request,omitempty"`
	MemoryNodeLimit   string `json:"memory_node_limit,omitempty"`
	JavaOpts          string `json:"java_opts,omitempty"`
	//DOCKERHOST
	ImageDockerhost         string `json:"image_dockerhost"`
	CpuDockerhostRequest    string `json:"cpu_dockerhost_request,omitempty"`
	CpuDockerhostLimit      string `json:"cpu_dockerhost_limit,omitempty"`
	MemoryDockerhostRequest string `json:"memory_dockerhost_request,omitempty"`
	MemoryDockerhostLimit   string `json:"memory_dockerhost_limit,omitempty"`
	PrivilegedDockerhost    bool   `json:"privileged_dockerhost,omitempty"`

	//CRAWLER
	ReplicasCrawler      int32  `json:"replicas_crawler,omitempty"`
	ImageCrawler         string `json:"image_crawler,omitempty"`
	GrpcAddresses        string `json:"grpc_adresses,omitempty"`
	CpuCrawlerRequest    string `json:"cpu_crawler_request,omitempty"`
	CpuCrawlerLimit      string `json:"cpu_crawler_limit,omitempty"`
	MemoryCrawlerRequest string `json:"memory_crawler_request,omitempty"`
	MemoryCrawlerLimit   string `json:"memory_crawler_limit,omitempty"`
	CrawlerServiceToken  string `json:"crawler_service_token,omitempty"`
	//AUTH
	ReplicasAuth           int32  `json:"replicas_auth,omitempty"`
	ImageAuth              string `json:"image_auth,omitempty"`
	CpuAuthRequest         string `json:"cpu_auth_request,omitempty"`
	CpuAuthLimit           string `json:"cpu_auth_limit,omitempty"`
	MemoryAuthRequest      string `json:"memory_auth_request,omitempty"`
	MemoryAuthLimit        string `json:"memory_auth_limit,omitempty"`
	MailEnabled            string `json:"mail_enabled,omitempty"`
	ActivateUserOnRegister string `json:"activate_user_enabled,omitempty"`

	//DATASERVICE
	ReplicasDs      int32  `json:"replicas_ds,omitempty"`
	ImageDs         string `json:"image_ds,omitempty"`
	CpuDsRequest    string `json:"cpu_ds_request,omitempty"`
	CpuDsLimit      string `json:"cpu_ds_limit,omitempty"`
	MemoryDsRequest string `json:"memory_ds_request,omitempty"`
	MemoryDsLimit   string `json:"memory_ds_limit,omitempty"`
	DsServiceToken  string `json:"dataservice_service_token,omitempty"`

	//FRONTEND
	ReplicasFrontend      int32  `json:"replicas_frontend,omitempty"`
	ImageFrontend         string `json:"image_frontend,omitempty"`
	CpuFrontendRequest    string `json:"cpu_frontend_request,omitempty"`
	CpuFrontendLimit      string `json:"cpu_frontend_limit,omitempty"`
	MemoryFrontendRequest string `json:"memory_frontend_request,omitempty"`
	MemoryFrontendLimit   string `json:"memory_frontend_limit,omitempty"`
	//AUTH_ADMIN
	ReplicasAuthAdmin int32  `json:"replicas_auth_admin,omitempty"`
	ImageAuthAdmin    string `json:"image_auth_admin,omitempty"`
	//Telegraf
	ReplicasTelegraf int32  `json:"replicas_telegraf,omitempty"`
	ImageTelegraf    string `json:"image_telegraf,omitempty"`
	//Nginx
	ReplicasNginx int32  `json:"replicas_nginx,omitempty"`
	ImageNginx    string `json:"image_nginx,omitempty"`
}

// WeMainnetStatus defines the observed state of WeMainnet
type WeMainnetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// WeMainnet is the Schema for the WeMainnets API
type WeMainnet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WeMainnetSpec   `json:"spec,omitempty"`
	Status WeMainnetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// WeMainnetList contains a list of WeMainnet
type WeMainnetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WeMainnet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WeMainnet{}, &WeMainnetList{})
}
