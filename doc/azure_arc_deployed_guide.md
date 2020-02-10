# Azure Arc API Gateway on K8S Cluster
## Architecture
![](https://i.imgur.com/SeS59ro.png)

## Running Azure API Management anywhere with Azure Arc for API Management
Azure Arc for API Management, a self-hosted Azure API management gateway (referred further to as self-hosted gateway as from now on), a self-hosted gateway that is a fully functional API Management gateway packed inside a docker container allowing you to run it anywhere you want.
![](https://i.imgur.com/Jbpwxih.png)

### Benefits
* Runs on any cloud, hybrid or on-premises – The self-hosted gateway, which is a full equivalent to a managed gateway, is packaged as a docker container meaning you can run it anywhere.

* Consistent API Experience –  Azure Arc for API Management enables companies to provide APIs consistently to their consumers. A win-win-win scenario for API consumers, operations and developer teams.

* Federated with API Management in the cloud – The self-hosted gateway will automatically pull down any new configuration deployed to the managed gateway in Azure and push its telemetry up. No inbound ports required, only outbound and no local storage required to maintain configuration data.

* Improved developer experience – It’s not the main reason why the self-hosted gateway was made but being able to clone an API, configure the API on a local self-hosted gateway and develop and test policies locally without interfering with other team members is definitely a nice bonus.

### Creating and running a self-hosted gateway
* To create a self-hosted gateway, navigate to your API Management instance and search for ‘Gateways’:
![](https://i.imgur.com/ajv4bXS.png)
* Press ‘Add’, and you’ll be prompted for a Name, Region and Description for your gateway, note that the Region has nothing to do with Azure Regions, its just a reference to a location that should be meaningful to where the self-hosted gateway is deployed.
![](https://i.imgur.com/M5c4mHe.png)
* To run an instance you’ll need to configure it via an environment file to tell it where the Azure API Management instance is located and register itself:
    * config.service.endpoint – Configuration endpoint of your Azure API Management instance.
    *  config.service.auth – Token used to register the gateway in you Azure API Management instance.
![](https://i.imgur.com/SF2WycK.png)
![](https://i.imgur.com/fYVbxGK.png)

### Using the self-hosted gateway
* Configuring the cat facts API on the gateway can easily be done from the gateway configuration page:
![](https://i.imgur.com/pml369t.png)
* Now let’s perform the same call as performed against the management API management instance bug against the IP address of the self-hosted gateway:
![](https://i.imgur.com/lotTF5M.png)

### Self-hosted gateway logging capabilities
Another powerful feature of Azure API Management is its out of the box metrics and analytics capabilities, let’s see how this extends to the self-hosted gateways.
* At container startup logging indicates the loading of configuration from the managed gateway in Azure and building of its internal routing table:
![](https://i.imgur.com/U8rzDth.png)
* Performing an update on the cat-facts API in the managed gateway in Azure triggers an immediate update on the local self-hosted gateway:
![](https://i.imgur.com/Y0BkJdF.png)

## Adding On-Premises Servers to Azure Arc
* Getting your on-prem servers to appear in the Azure Arc portal is pretty straight forward. First, we need to make sure we have a few things checked off before we dive in.
* Required Resource Providers -We need to register two resource providers to use Azure Arc for Servers.
    * Microsoft.HybridComputer
    * Microsoft.GuestConfiguration
* This can be done in the portal, through PowerShell or Azure CLI. We’ll be using PowerShell in this example.
![](https://i.imgur.com/ejOXiY9.png)

* We can also verify in the portal as well.
    * Click on Subscriptions.
    * Choose your subscription.
    * Under settings, select Resource providers.
    * Using the filter by name option and locate GuestConfiguration and HybridComputer
* GuestConfiguration resource provider seen below.
![](https://i.imgur.com/YBcgIGl.png)
* HybridCompute resource provider seen below.
![](https://i.imgur.com/NdpemLy.png)

### Adding Servers
Once we have that complete, we can add our servers by completing the following within the Azure portal.
* Type Azure Arc in the Search resources box at the top of the port and press Enter (see below).
![](https://i.imgur.com/L8thhdO.png)
* Select Manage Servers while in the Azure Arc portal.
* In the Machines screen. Click + Add to add a server.
![](https://i.imgur.com/1YD1P4r.png)
* Choose Generate Script in the Add machines using interactive script.
* Select your subscription, resource group (or create new) and region.
![](https://i.imgur.com/oNhA8Q8.png)
* Choose the operating system your on-premises workload is running, either Windows or Linux. We’ll choose Windows.
* Click Review and generate script
![](https://i.imgur.com/g5Std93.png)
* At this point, you can either download or copy the script.
* After you have either downloaded or copied the script. You’ll need to run these commands on the server you want managed with Azure Arc. The script downloads a lightweight agent and installs it on the server which in turn associates with your subscription. Your server should appear in the Azure Arc portal after several minutes.
* If you’re on-boarding a Linux server, you’ll have to copy the commands from the portal and execute on the server. They do not provide a script to download. The Linux commands perform the same process as the Windows commands by downloading and installing a lightweight agent on the server.
![](https://i.imgur.com/q11qKYp.png)

## References
* https://www.codit.eu/blog/running-azure-api-management-anywhere-with-azure-arc-for-api-management/
* https://docs.microsoft.com/zh-tw/azure/api-management/api-management-howto-deploy-self-hosted-gateway-to-k8s
* https://www.petri.com/getting-started-with-azure-arc-servers



