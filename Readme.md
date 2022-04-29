Variables:

=========NODE=========

required:

	image:
	replicas: 3
	storage: "2Gi"

optional:

	clean_state: "false" - default
	cpu_node_request
	cpu_node_limit
	memory_node_request
	memory_node_limit
	java_opts: "-Dwe.check-resources=false -Xmx3g" - default

========CRAWLER========== 

required:

	replicas_crawler
	image_crawler

optional:

	grpc_adresses
	cpu_crawler_request
	cpu_crawler_limit
	memory_crawler_request
	memory_crawler_limit
	crawler_service_token

=========AUTH=========== 

required:

	replicas_auth
	image_auth

optional:

	activate_user_on_register
	cpu_auth_request
	cpu_auth_limit
	memory_auth_request
	memory_auth_limit
	mail_enabled

========DATASERVICE========

required:

	replicas_ds: 1
	image_ds

optional:

	cpu_ds_request
	cpu_ds_limit
	memory_ds_request
	memory_ds_limit
	dataservice_service_token

========FRONTEND========= 

required:

	replicas_frontend: 1
	image_frontend

optional:

	cpu_frontend_request
	cpu_frontend_limit
	memory_frontend_request
	memory_frontend_limit
	
========AUTH-ADMIN=========

	replicas_auth_admin: 1
	image_auth_admin
