General


The Kubernetes operator is designed to automate application deployment. You only need to install it and fill in the parameters in the configuration file. All services will be automatically deployed. Allows you to use a set of services as one application, by simple configuration. Using the settings, you can choose to install with a client, only a node, and any variations. The operator independently monitors that all services are deployed and running in the cluster.

Main settings


To create a WeMainnet service after installing the operator, you need to create a service with kind: WeMainnet. The minimum configuration looks like this:

```
apiVersion: apps.wavesenterprise.com/v1
kind: We Mainnet
metadata:
  name: mainnet
  postgres_enabled: "false"
#NODE
  image: "wavesenterprise/node:v1.8.1"
  replicas: 1
  storage: "20Gi"
#DOCKERHOST
  image_dockerhost: "docker:19.03.1-dind"
  privileged_dockerhost: true
```

In this case, the following services will be deployed:

node
dockerhost


First you need to add configs for nodes and wallets to the namespace:

```
---
apiVersion: v1
kind: Secret
metadata:
  name: node-wallet
type: Opaque
data:
  node-0-keystore.dat: /FILL IN YOUR KEYSTORE IN BASE64/

---
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-config
data:
  node-0.conf: |
     /FILL IN YOUR NODE CONFIG/
```



Additionally, you can deploy the telegraf agent

spec:
  telegraf_enabled: "false"


Telegraf is used as an agent for building node metrics, it must be specified in the node configuration in the metrics.uri section

telegraf requires configuration

```
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: telegraf-config
data:
  telegraf.conf: |
    /FILL IN YOUR TELEGRAF CONFIG/
```

The postgres_enabled: "false" option can be added to the postgresql database deployment, databases for client services will be automatically created.



When the option client_enabled: "true" is enabled, client services will be deployed:

Crawler
Auth Service
data service
frontend
The following parameters in the manifest become mandatory:

```
#CRAWLER
  replicas_crawler: 1
  image_crawler: "wavesenterprise/crawler2:v1.8.0"
#AUTH
  replicas_auth: 1
  image_auth: "wavesenterprise/auth-service:v1.8.0"
#DATASERVICE
  replicas_ds: 1
  image_ds: "wavesenterprise/data-service:v1.8.0"
#FRONTEND
  replicas_frontend: 1
  image_frontend: "wavesenterprise/frontend-app:v1.8.0"
#AUTH ADMIN
  replicas_auth_admin: 1
  image_auth_admin: "wavesenterprise/auth-service-admin:v1.8.0"
```

In some cases, the web interface of the authorization service may be useful, it can also be additionally expanded with the parameter

spec:
  authadmin_enabled: "false"




Extra options:
The parameters are set by default, they can be added to the manifest if you want to change

node:
```
clean_state: "false"

cpu_node_request: "2"

cpu_node_limit: "2"

memory_node_request: "4Gi"

memory_node_request: "4Gi"

java_opts: "-Dwe.check-resources=false -Xmx3g"
```



dockerhost:
```
cpu_dockerhost_request: "1"

cpu_dockerhost_limit: "2"

memory_dockerhost_request: "1Gi"

memory_dockerhost_limit: "2Gi"
```



Crawler
```
grpc_addresses: "node-0:6865,node-1:6865,node-2:6865"

cpu_crawler_request: "1"

cpu_crawler_limit: "1"

memory_crawler_request: ""2Gi"

memory_crawler_limit: "2Gi"

crawler_service_token: ""
```



Auth Service:
```
cpu_auth_request: "100m"

cpu_auth_limit: "1"

memory_auth_request: "32Mi"

memory_auth_limit: "1Gi"

mail_enabled: "false"
```


Data service:
```
cpu_ds_request: "1"

cpu_ds_limit: "1"

memory_ds_request: "1Gi"

memory_ds_limit: "1Gi"

dataservice_service_token: ""
```


frontend:
```
cpu_frontend_request: "1"

cpu_frontend_limit: "1"

memory_frontend_request: "1Gi"

memory_frontend_limit: "1Gi"
```


Configuration files required for services


Installation without a client:

Through this configuration, you can declare Env variables in the node container, the config is required

```
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: vars-node
data:
  WE_NODE_OWNER_PASSWORD: "/PASSWORD/"
  WE_NODE_OWNER_PASSWORD_EMPTY: "false"
```


This configuration contains the config for the node, you need to fill it with your own parameters

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-config
data:
  node-0.conf: |
    node {
      waves-crypto="yes"
      owner-address = "/ADDRESS/"
      directory="/node"
      license.file = "/opt/licenses/node-0.license"
      data-directory = "/node/data"
      wallet {
        file = "/opt/wallets/node-0-keystore.dat"
        password= "/PASSWORD/"
      }
      blockchain {
        type="MAINNET"
      }
      logging-level="DEBUG"
      network {
        bind-address="0.0.0.0"
        port=6864
        known-peers = [
          "cloud.wemeadow.com:6864"
          "pool.wemeadow.com:6864"
          "node-0.wavesenterprise.com:6864"
          "node-1.wavesenterprise.com:6864"
          "node-2.wavesenterprise.com:6864"
        ]
        node-name="MAINNET_NODE_2"
        peers-data-residence-time="2h"
        break-idle-connections-timeout = "3m"
        declared-address = "0.0.0.0:6864"
      }
      api {
        auth {
          type = "api-key"
          api-key-hash = "/KEYHASH/"
          privacy-api-key-hash = "/PRIVACYKEYHASH/"
        }
      }
      miner {
        enable = "yes"
        quorum = 2
        interval-after-last-block-then-generation-is-allowed = "35d"
        micro-block-interval = "5s"
        min-micro-block-age = "3s"
        max-transactions-in-micro-block = 500
        minimal-block-generation-offset = "200ms"
      }
      scheduler-service {
        enable = "no"
      }
      privacy {
        crawling-parallelism = 100
        storage {
          vendor = "none"
        }
      }
      docker-engine {
        enable = "yes"
        use-node-docker-host = "yes"
        default-registry-domain = "registry.wavesenterprise.com/waves-enterprise-public"
        docker-host = "tcp://dockerhost:2375"
        execution-limits {
          timeout = "30s"
          memory = 512
          memory-swap = 0
        }
        reuse-containers = "yes"
        remove-container-after = "10m"
        remote-registries = []
        check-registry-auth-on-startup = "yes"
        contract-execution-messages-cache {
          expire-after = "60m"
          max-buffer-size = 10
          max-buffer-time = "100ms"
        }
      }
    }
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: node-license
data:
  node-0.license: |

    {"license":"....}

```

For the authorization service, you need to generate keys and add them to the auth-service-keys secret

Generating a new key pair:

```
ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key
 # Don't add passphrase
 openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
 cat jwtRS256.key | base64
 cat jwtRS256.key.pub | base64
```

```
---
apiVersion: v1
kind: Secret
metadata:
  name: auth-service-keys
type: Opaque
data:
  jwtRS256.key: LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS0FJQkFBS0NBZ0VBa3FmcW1PSUNTQ203ZFZBUjlnOEhZSXFiNUhqd05KNDdPK1JLVDZMV3lSSzBBN084Ckp6R2FaaENRdjFOMmtqeXloSHdtMHN1RUZxRzhkMHRQb3Rvb0l6ZXF4eFFKWHd1ZVlHU1hseWZHNHJMcEJqUVUKQktZZTVtcytLR05ZeFZ0SFNFUjROQVpnNVNBSDEwcjR3ck9MKy81U0xaWk9WTkliajU0YlY4S1BXNFJIU3JNeQp5YmcwTDNyN1FJekc5SEpzaU1WcmIzTjVPTnZDalNvb2pVclBkNGdoSm04NjNFSmYzRFdZbXY1bmx4UjlXbHI3CnZSbUxaS3JsS2N0cWpJaHgxZjVnT3BFcjVVclFwRGZDMXBFaDMrVUs0K1FUcHN6YUdaRXBnNndMRmJpZHBNYWgKUVU3cjA5WWZLYnhwSm1udldSK2lYaWdIWnBJWE5kMzZBanRlTkhTMnFIdmJTVTVKTmZDRDdndTA0dzlwL3BVTgozTGdDYVdwZDFKd1doUFFQQVZOWGVDaWdDM044aDBqWGRUcThtMTNDNzIwMDNlR0REZ0dNSW5xR0xKdVRJY0tFCjF6UEQzSUZZczVVZ3ZNMElaQXlodHljellEV3hCNFZPSFNLQ21ZUThGQWF6b05aTTBBd0kvUE5PYm8xTU1ienAKc2tsWFJSRXZFVDJaUGtBMHAwU2xjQ1ArcEVRdVZ2dVZ3dG9aMzV0OU9DS0t6djhlQmJYRDdqbWJpdU80enZXRQpHNUdoZXg2YUhjMGZ2cDcvUTJRUExZU2s0UnQ3REFMc1hjZVBxeWJnbjEyTE1lMGo1OVRielFtdW1LdUd4bS9KClZ1anUyYkQ0dWVyUXBLaEg3M1dkdHFMSFRiZ3JBRmdiWms3V1RhRWVCbWlvVlRsenQ3KzN5V0dpblJVQ0F3RUEKQVFLQ0FnQXU3aVVZZjFxVGxTY3p4MGU0SWQ4T2VjeWhORUpKMUVqSVJhbXlDajRKWWo3UTRIZFpZM294SnlQcQoxZDZmdFdTN0dLK2p6UlNiMlczaUR1dVJCWmJLamtuTUl4Rk1wUDh0Z1lNeEQ3MkpWZzlUdU42ZkRqbmRLbnhkCm1FMFQrcjI0MXBCUXRhblVLSWZaMFZnQmxrczVmSXozb1oyM2J2VDY1SEdEaC9Nd0tnaFdVem54YTB6bjFNY0sKUlFKMXZ4Z1VQSGpBMVliNU52bnZDb3FuakVVUHp6UXNoSE9sZ0dnRW8vSU54MU9HK1R1VDZvR2NaY3hCanA1KwozV1ZmUzFxQ1RFQ20vVVc2dmxJOStzb1N0NmJMYXhRdEVSTm8wUzBKK1hYN2VOYWpRTXpScWU1NFk3VDd1UTJICjRZOFVoOW5iLzArS0tlMFVXYk9ydlRqLzlkZVRIYktvUGlMcVVoS2ZLK0pwWWV2YnhLT242eVk5UFhRQndLN0gKQ3EyekxxNzFkUDJDc3F0WGplOWxVYjEwWTAzdHdoWG5qSzQ1Q2RvMUl0aXAzQlNERmVINHFHTXNHVVgyTVpUdgpwUmpQWHEraG16c3pKUnFoUWs2K3NMVU1OU1pkTVdtZzFpNUNpdkZQM2kybnczbW1ma3RXL3hMTnRMUzlPZkw0Cmd0UHRQSVlXUTBSdDdCRDFDak1RM3RFeGxldTZmd1lVbmlQTmtySmZFMEhXaGdnWUtpUCtjRXN0cWg3VzRweDQKS2VYY3EzR0hubDRKNm9FRTZpd0lkdzdoL2tEcHdDdStzeC9Ld2svNkJ3R0xIQisvTzIrK2dPZGxFRHJvK1EySgpjVzZGQU5rbmowTGp5RVNJZExrVUsrTEd2TmNkSjZVQ0xMcjhKZElIVTFuL0U4NEpFUUtDQVFFQXd5ek1vNVZBCkhxKzMvWWhEby9kVVY4eFZqV01tYjdUWXBMOEFiZ3ZQTjE3RjBMOUoySDZoamNGbm1XT3o5VU1sZ0IxN3Y0Mk0KcFRLdnREWHZlSjlVTXlHT2dVZ0tTNElqTkdudy9KSGZ6NWZjVzhicmpXVi9rVUIxQlhzMVovbEY2Q0RSYmNOcwphT0xNY0UxK3ZmOGZCRmNvaldhRHZ4M1VtL0M5ekp3ZXRuUzBlWGtmR1lrV2VFdWU5QUkyVjgyZWFjM3JGeXpwCkpXUFdld3k4MHI4ZTA5TE1WK21CR2w4ZFhxMXVsREpULzBpUnppb1c0U2wrcERpZytvbXV5K2pjb1dGZnc4bisKQ1R1czlwcG5ZYjZJZk1mOXA5TzQxUHVONzQ5S01Ca3E0NzB3a0VVMEp6SlBldmJNYkpUYmFmTzd3VTBhNndZNwpOUlBpU0ZJREN6U2Nvd0tDQVFFQXdGdzdERHhjMElXYXJYdWdNa0orS0toTU9YR3RabEQrTzlMNnZFUVBGZmVFCnAxTXhNUTA0a1FkbWtYekRFWmVHWG9GQVIrQmdCMzRJc2d0OFBTWVdKUmtxMVhwMmh0N2JDNmxjcTAxVTZrL0QKZUE4YjNtSnQ0eTZ4UkZzU2ttc0dkMEdvMSsrb3EzbWF1L0NMVkRNMkRFZUpxemZWQzhnS0ZYRkZkSE04SlN5bQo4dUtyd1dENzllQWtjNDUzUjdCOXpicVZtVURBYXlLQmxLVFA4ZGJ0T3loQzZ0UVpBblY1enJvZFdPSm1LaWtDCldMNUVld1J5K3JDdnBWTE5TbXhEV2V6VWZPcVZ3MEhmUi9pZWVHVHFBb1UrRGo5Ri82Ti9kcm1IOWN3L0NWeVoKWkVRY2wwL1lZZ0F3N2dkUThPdzJsVHNZcnU4d1Vqb1Rtemh2QTFBQzV3S0NBUUFtQks0QUYyeWNEYUtMY21XcQpwTno3RlVSOC9CbGFuU0d1UmI1eHNUODJDL0lBamFKMjE0UGt0dzNWSlVUQ3U4ZXNReEg5NkRiRFh6STJxbUx4ClhpZnFwZGk2ZWl2M05XeGlJMWpiK2haY3U3b2k3b2Fuem1PaENhdEIzQlExSXF0cFlpc3BkRzNEcUpvbUxoSkkKTkUvNGFubnR3VkJjaEJVTUkwTDFmbHZGTXNxTTl2a0Y0bHhNSm43YURTeEV3anJmWlVzc0FvV1AwUGpRazFTYwp3TG5palNkYzRKRlRiNytxMTZHNG9HMFlSeXlQdWtjbXFReVFOSysyM2ViOHRXbDB6aUQzWkh0bGxRaEdLU0dHCk9yVWZpVjF4dVo1QmJwYmhXVW9jUUdySVhldjl6bDB3WFc1NkIyVWVxWWhzQlJ4SHRSdFBPTEdEejFHK3dLcnoKSGRaOUFvSUJBUUNkSi8wWi93cjVWZDVNVkE5S1lLYS93dGdicW5NM2YzNW1FL1hEOEhxK3dLMUJJeWV5WXBIUApjMU5xRTVzdmVUTlBiSnUrM1dLM1hGSHdYSS9SU1plWUVacThOTVEzWmtWaG5xbldUbVRNMWdQbHg3cEdFdmFpClFCaVZ0eTVTTDF4bC9GL2NvN0dTL3RQYkxpZzJ6MndkMWloMG1UWFczVVRYeGVZdndLSG40VFk0ZzlZOU5HWkYKdUMwdnQ3cGQrS1NmZXd3VDNDSVlwV1Ztc3N3dFVpSVpUY2gySUhpYVdLMytwbkdwbDdaT1JaamtOZmF1NXJDbApmY3JTNy9aSEVuSm9PcVJUdGpoTEFUdFJpcDYxMEFTYnNJNUZoNDVCMENzb0xXWVYvQnVZSTI0eXk2N3NORkNkCnFIaFJUK3JpR3FweGU0bXNDa0RaUFJlZG5ocWNnemNMQW9JQkFHWjdzRjRzd0JUbGNqOHcrUGg3cXRQaVVGSEQKNzhmU3JkQ1RtTUF2eFFKbHdTVjZENjFndmx6NUZoamVOdmVnd0JmdzRUTlhxYXlyWjVXZHFGQlZXRWhzTlRNbQorOVBaR1pUcWdML1pPckZycTF4WUVGeVlVYStaOUVXcWJxa1hOekFhWWJyY0RIRC8wdEJJRkk0RlFLRjhER1pICnVET1lRSzN1UkdjWFFVem9GTDZzNTVCdFliZE03NzRieGt1d1l2MVp5U0RabGVyWXgvbU1YWXYzWlpuSC9ZYjgKakRKWmRLZjdQeEQxd25KR2Q3L09nZFRIdSs4SmxDMXMvbHF4Y1kvR1VqOVozdksrVnBvOFhkZ2E2TEQvY1NLaQo2bXMrcFRZRUh2b043c1k1Zi9vTTRMK2k3RVZCWWxqTDRQbWZza2Y1dnJ5Qm1jZ0ZPZGYyNWw5cVllND0KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0=
  jwtRS256.key.pub: LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUlJQ0lqQU5CZ2txaGtpRzl3MEJBUUVGQUFPQ0FnOEFNSUlDQ2dLQ0FnRUFrcWZxbU9JQ1NDbTdkVkFSOWc4SApZSXFiNUhqd05KNDdPK1JLVDZMV3lSSzBBN084SnpHYVpoQ1F2MU4ya2p5eWhId20wc3VFRnFHOGQwdFBvdG9vCkl6ZXF4eFFKWHd1ZVlHU1hseWZHNHJMcEJqUVVCS1llNW1zK0tHTll4VnRIU0VSNE5BWmc1U0FIMTByNHdyT0wKKy81U0xaWk9WTkliajU0YlY4S1BXNFJIU3JNeXliZzBMM3I3UUl6RzlISnNpTVZyYjNONU9OdkNqU29valVyUApkNGdoSm04NjNFSmYzRFdZbXY1bmx4UjlXbHI3dlJtTFpLcmxLY3RxakloeDFmNWdPcEVyNVVyUXBEZkMxcEVoCjMrVUs0K1FUcHN6YUdaRXBnNndMRmJpZHBNYWhRVTdyMDlZZktieHBKbW52V1IraVhpZ0hacElYTmQzNkFqdGUKTkhTMnFIdmJTVTVKTmZDRDdndTA0dzlwL3BVTjNMZ0NhV3BkMUp3V2hQUVBBVk5YZUNpZ0MzTjhoMGpYZFRxOAptMTNDNzIwMDNlR0REZ0dNSW5xR0xKdVRJY0tFMXpQRDNJRllzNVVndk0wSVpBeWh0eWN6WURXeEI0Vk9IU0tDCm1ZUThGQWF6b05aTTBBd0kvUE5PYm8xTU1ienBza2xYUlJFdkVUMlpQa0EwcDBTbGNDUCtwRVF1VnZ1Vnd0b1oKMzV0OU9DS0t6djhlQmJYRDdqbWJpdU80enZXRUc1R2hleDZhSGMwZnZwNy9RMlFQTFlTazRSdDdEQUxzWGNlUApxeWJnbjEyTE1lMGo1OVRielFtdW1LdUd4bS9KVnVqdTJiRDR1ZXJRcEtoSDczV2R0cUxIVGJnckFGZ2JaazdXClRhRWVCbWlvVlRsenQ3KzN5V0dpblJVQ0F3RUFBUT09Ci0tLS0tRU5EIFBVQkxJQyBLRVktLS0tLQ==

---
apiVersion: v1
kind: Secret
metadata:
  name: node-wallet
type: Opaque
data:
  node-0-keystore.dat: /WALLET BASE64/
  
```

Using the auth-service-tokens secret, you can set your application access tokens. They will need to be set via a parameter in the crawler and dataservice.

```
---
apiVersion: v1
kind: Secret
metadata:
  name: auth-service-tokens
type: Opaque
data:
  tokens.json: ICAgIHsKICAgICAgImFjY2VzcyI6ICJNaDdLMGlNajFqZTNwckMza2tqaFpFM0ZPWDVpblBSYlh2T0Foc1BSIgogICAgfQ==
```

Using the pg-user secret, we set parameters for accessing services to the database

```
---
apiVersion: v1
kind: Secret
metadata:
  name: pg-user
type: Opaque
data:
  user: cG9zdGdyZXM=
  password_admin: VE1MdDFoUE9GVG8xMG9Gaw==
  ssl: ZmFsc2UK

```

Using a mail secret, you can set the parameters for sending emails via smtp

```
---
apiVersion: v1
kind: Secret
metadata:
  name: mail
type: Opaque
data:
  mail-host: ZXhhbXBsZQ==
  mail-user: ZXhhbXBsZQ==
  mail-password: ZXhhbXBsZQ==
  mail-from: ZXhhbXBsZQ==
  mail-port: ZXhhbXBsZQ==
  mail-salt: ZXhhbXBsZQ==
  ```

