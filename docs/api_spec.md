API Spec
========


# Model

# Command Model

* Starts with prefix "/": URI Model for command logging
* Starts without prefix "/": Command for sync model 

## Sync Model
| Name | Type | Default Value | Req. Required |
|------|------|---------------|---------------|
|timestamp|Timestamp|-|V|
|version|String|-|-|
|command|String|-|V|

### Available Commands
| Command | Description |
|---------|-------------|
| HELO    | Indicate HELLO from client |
| IREQ    | Send identity request to client |
| IRESP	  | Send identity reply from client |
| TREQ    | Send timestamp sync request with server timestamp and version to client |
| TRESP   | Send timestamp sync with client timestamp and version from client |
| TOK     | Send sync timestamp complete to server |
| SREQ    | Send commands fetch that out of sync (if required) from server |
| SREP    | Send commands by server requested |
| SOK     | Send command synced from client, indicate command is sync between client and server |  


## Data Model
| Name | Type | Default Value | Req. Required | Available options |
|------|------|---------------|---------------|-------------------|
|id|String|-|V|Primary key for data model|
|name|String|-|V|-|
|type|String|-|V|Model; Metadata; YAML; Configuration|
|extension|String|-|V|.tar.gz; .zip; .7z; .yaml; |
|external|Boolean| False | V |True; False |
|location|String|-| - | -| 

## Container Model
| Name | Type | Default Value | Req. Required | Available options |
|------|------|---------------|---------------|-------------------|
|id|String|-|V|Primary key for container model|
|name|String|-|V|-|
|status|String|-|V|-|



# API List

## Sync
| Method | URI | Parameter | Description |Return| Comment |
|--------|-----|-----------|-------------|------|---------|
|GET |/v1/sync|(None)|Get current sync status|[Sync]|(None)|
|POST|/v1/sync|command:Base64Encoded URI String|Binding URI to new version|[Sync]|(none)|

## Data (Model/Metadata)
| Method | URI | Parameter | Description |Return| Comment |
|--------|-----|-----------|-------------|------|---------|
|GET|/v1/data|-|-|-|List function not applicable|
|POST|/v1/data|[Data]|Data model about created data|-|-|
|GET|/v1/data/[Id]|-|Get data model by Id|-|-|
|PUT|/v1/data/[Id]|[Data]|Update data model by Id|-|-|
|DELETE|/v1/data/[Id]|-|Delete data model by Id|-|-|

## Container
| Method | URI | Parameter | Description |Return| Comment |
|--------|-----|-----------|-------------|------|---------|
|GET|/v1/container|offset: integer|List available containers|[Container]|-|
|POST|/v1/container|[Container]|Add new container into system|-|-|
|GET|/v1/container/[Id]|-|Get container info by Id|-|-|
|PUT|/v1/container/[Id]|[Container]|Update container info by given model|-|-|
|DELETE|/v1/container/[Id]|-|Delete container by Id|-|-|

## System (Kubernetes & Machine) 
| Method | URI | Parameter | Description |Return| Comment |
|--------|-----|-----------|-------------|------|---------|
|GET|/v1/system/info|-|-|-|-|
|[METHOD]|/v1/system/[RES]/[ID]|-|Run specified resource command.|-|See below|

> The system mapping commands are available for CMS to control remote cluster by local kubernetes control commands via transparent mode. The METHOD, RES, ID, and along with command parameters will transfer to equivalent resource and operations in the cluster.
> 
> For example, the command POST /v1/system/deployment with deployment YAML as model will trigger remote cluster create a new Deployment, same effort as "kubectl create deploy -f 'Deploy.yaml'".

