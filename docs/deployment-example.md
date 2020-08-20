 Deployment Example
 ==================
 [TOC]

 To deploy an application on the itrix-edge trigger by agent, users must do the following steps:

1. Register the deployment
2. Trigger the deployment

# Register the deployment

## Prepare the deployment-related JSON

 There are at least two JSON definitions need to register, one is `Deployment`, another one is `Service`.

 Users should use standard kubernetes deployment and service JSON to complete the registeration process. The YAML template can be find here, or use convenience template generation provided by `kubectl` command line:

```=shell
# On Ansible-host, with inside the ITRIX-EDGE folder.
$ ./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf create deployment [application-name] --image=[application-image] --dry-run=true -o json > application-name.json
```
 The command will generate a valid `Deployment` JSON model with the following content:
```=json
{
    "kind": "Deployment",
    "apiVersion": "apps/v1",
    "metadata": {
        "name": "application-name",
        "creationTimestamp": null,
        "labels": {
            "app": "application-name"
        }
    },
    "spec": {
        "replicas": 1,
        "selector": {
            "matchLabels": {
                "app": "application-name"
            }
        },
        "template": {
            "metadata": {
                "creationTimestamp": null,
                "labels": {
                    "app": "application-name"
                }
            },
            "spec": {
                "containers": [
                    {
                        "name": "application-image",
                        "image": "application-image",
                        "resources": {}
                    }
                ]
            }
        },
        "strategy": {}
    },
    "status": {}
}
```
 To generate service template yaml, use the following command instead:
```=shell
# On Ansible-host, with inside the ITRIX-EDGE folder.

# Used to expose application to the internet:
$ ./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf create service loadbalancer [service-name] [--tcp=port:targetPort[,port2:targetPort2,...]] --dry-run=true -o json > service-name.json

# Used to NOT expose application:
$ ./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf create service clusterip [service-name] [--tcp=port:targetPort[,port2:targetPort2,...]] --dry-run=true -o json > service-name.json
```
 The command will generate a valid `Service` JSON model with the following content:
```=json
{
    "kind": "Service",
    "apiVersion": "v1",
    "metadata": {
        "name": "service-name",
        "creationTimestamp": null,
        "labels": {
            "app": "service-name"
        }
    },
    "spec": {
        "ports": [
            {
                "name": "port-targetPort",
                "protocol": "TCP",
                "port": port,
                "targetPort": targetPort
            },
            {
                "name": "port2-targetPort2",
                "protocol": "TCP",
                "port": port2,
                "targetPort": targetPort2
            }
        ],
        "selector": {
            "app": "service-name"
        },
        "type": "LoadBalancer"  # Or "ClusterIP"
    },
    "status": {
        "loadBalancer": {}
    }
}
```

## Register deployment and service as the template

```=shell
# On Ansible-host, with inside the ITRIX-EDGE folder.
$ export AGENT_IP=`./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf ge t svc -n edge-client-agent edge-agent -o go-template='{{.spec.clusterIP}}'`
$ curl -X POST -k -v "http://$AGENT_IP:9000/v1/deplymentTemplate" -d '{JSON_Data}'
```
### JSON data options for register deployment

 | Options       | Default Value | Comment |
 | ------------- | ------------- | ------- |
 | namespace       |default | running namespace |
 | options       |[{"key":"k", "value": "v"}] | template variable definition |
 | deployment_template    |`Deployment` default model | `Deployment` model (Note 1) |
 | service_template       |`Service` default model | `Service` model (Note 2)|

 Note:
1. Fill in `Deployment` model in JSON format.
2. Fill in `Service` model in JSON format.

Sample JSON postdata for the registration request:
```=json
{
    "namespace": "default",
    "options": [
        {
            "key": "template_key",
            "value": "template_value"
        }
    ],
    "deployment_template": {
        "kind": "Deployment",
        "apiVersion": "apps/v1",
        "metadata": {
            "name": "application-name",
            "creationTimestamp": null,
            "labels": {
                "app": "application-name"
            }
        },
        "spec": {
            "replicas": 1,
            "selector": {
                "matchLabels": {
                    "app": "application-name"
                }
            },
            "template": {
                "metadata": {
                    "creationTimestamp": null,
                    "labels": {
                        "app": "application-name"
                    }
                },
                "spec": {
                    "containers": [
                        {
                            "name": "application-image",
                            "image": "application-image",
                            "resources": {}
                        }
                    ]
                }
            },
            "strategy": {}
        },
        "status": {}
    },
    "service_template": {
        "kind": "Service",
        "apiVersion": "v1",
        "metadata": {
            "name": "service-name",
            "creationTimestamp": null,
            "labels": {
                "app": "service-name"
            }
        },
        "spec": {
            "ports": [
                {
                    "name": "port-targetPort",
                    "protocol": "TCP",
                    "port": "port",
                    "targetPort": "targetPort"
                },
                {
                    "name": "port2-targetPort2",
                    "protocol": "TCP",
                    "port": "port2",
                    "targetPort": "targetPort2"
                }
            ],
            "selector": {
                "app": "service-name"
            },
            "type": "LoadBalancer"  # Or "ClusterIP"
        },
        "status": {
            "loadBalancer": {}
        }
    }
}
```
Below is the full registration example:

```=shell
# On Ansible-host, with inside the ITRIX-EDGE folder.

# Below for agent IP INSIDE the cluster
$ export AGENT_IP=`./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf get svc -n edge-client-agent edge-agent -o go-template='{{.spec.clusterIP}}'`
$ curl -X POST -k -v "http://$AGENT_IP:9000/v1/deploymentTemplate" -d '{"namespace":"default","options":[{"key":"template_key","value":"template_value"}],"deployment_template":{"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"application-name","creationTimestamp":null,"labels":{"app":"application-name"}},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"application-name"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"application-name"}},"spec":{"containers":[{"name":"application-image","image":"application-image","resources":{}}]}},"strategy":{}},"status":{}},"service_template":{"kind":"Service","apiVersion":"v1","metadata":{"name":"service-name","creationTimestamp":null,"labels":{"app":"service-name"}},"spec":{"ports":[{"name":"port-targetPort","protocol":"TCP","port":80,"targetPort":80},{"name":"port2-targetPort2","protocol":"TCP","port":8080,"targetPort":8080}],"selector":{"app":"service-name"},"type":"ClusterIP"},"status":{"loadBalancer":{}}}}'

< HTTP/1.1 200 OK
< Access-Control-Allow-Credentials: true
< Access-Control-Allow-Headers: X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token
< Access-Control-Allow-Methods: POST, GET, OPTIONS, PUT, DELETE, UPDATE
< Access-Control-Allow-Origin: http://localhost
< Access-Control-Expose-Headers: Content-Length
< Access-Control-Max-Age: 86400
< Content-Type: application/json; charset=utf-8
< X-Request-Id: d65453d5-c97d-48b4-b35a-b5ada1ee6f0f
< Date: Wed, 19 Aug 2020 09:38:35 GMT
< Content-Length: 1367
<
{"data":{"id":3,"created_at":"2020-08-19T09:38:35.005466439Z","updated_at":"2020-08-19T09:38:35.005466439Z","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"namespace":"default","options":[{"key":"template_key","value":"template_value"}],"deployment_template":{"kind":"Deployment","apiVersion":"apps/v1","metadata":{"name":"application-name","creationTimestamp":null,"labels":{"app":"application-name"}},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"application-name"}},"template":{"metadata":{"creationTimestamp":null,"labels":{"app":"application-name"}},"spec":{"containers":[{"name":"application-image","image":"application-image","resources":{}}]}},"strategy":{}},"status":{}},"service_template":{"kind":"Service","apiVersion":"v1","metadata":{"name":"service-name","creationTimestamp":null,"labels":{"app":"service-name"}},"spec":{"ports":[{"name":"port-targetPort","protocol":"TCP","port":80,"targetPort":80},{"name":"port2-targetPort2","protocol":"TCP","port":8080,"targetPort":8080}],"selector":{"app":"service-name"},"type":"ClusterIP"},"status":{"loadBalancer":{}}},"hooks":[{"id":3,"created_at":"2020-08-19T09:38:35.007220672Z","updated_at":"2020-08-19T09:38:35.007220672Z","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"name":"default.application-name","key":"3c4ef1493e7de67ffdce0ad7100933e3","deplyoment_option_id":3}]}}
* Connection #0 to host 10.233.41.19 left intact
```

 Where the last line with the `"key":"3c4ef1493e7de67ffdce0ad7100933e3"` is the deployment trigger key generated by the agent.

# Trigger the deployment

Its simple to start a registered deployment by hit the trigger with the given key `3c4ef1493e7de67ffdce0ad7100933e3`:

```=shell
# On Ansible-host, with inside the ITRIX-EDGE folder.

# Below for agent IP INSIDE the cluster
$ export AGENT_IP=`./kubespray/inventory/edge/artifacts/kubectl --kubeconfig=./kubespray/inventory/edge/artifacts/admin.conf get svc -n edge-client-agent edge-agent -o go-template='{{.spec.clusterIP}}'`

# Trigger to start the deployment
$ curl -X GET -k -v http://$AGENT_IP:9000/v1/key/3c4ef1493e7de67ffdce0ad7100933e3

# You could use POST to bring variables to the edge agent:
$ curl -X POST -k -v http://$AGENT_IP:9000/v1/key/3c4ef1493e7de67ffdce0ad7100933e3 -d '[POST_DATA]'
```

After successfully trigger the deployment, use usual Kubernetes management tools or `kubectl` to manage the deployment.