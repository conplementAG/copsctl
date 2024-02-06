# Build Agent Pool

You can manage your project specific build agent pool based on azure virtual machine sets. The build agents are 
linux based with docker support. Teams creating build agent pools are responsible for the created pools by their own.

## Preconditions
Required tools
- Terraform >= 1.7
- Azure Cli >= 2.56.0
- Sops >= 3.8.1

The user running the commands needs owner permissions on azure subscription infrastructure will be created.  

## Infrastructure and Cost
The build infrastructure created consists of a virtual machine scaleset and an egress loadbalancer with static 
IP address. VMSS is by default configured to scale down to 0 if not in use to reduce cost. The total resources cost 
is depended on the amount of builds running and configuration. Expect a monthly cost between 50€ und 100€.

## Setup
1. Create a sops encrypted configuration 
2. Run create command
3. Ask cops team to add your build agent IP address to cluster firewalling

### Configuration
[Config Template](./configuration/.corebuild.template.yaml)

### Create
```
copsctl build create -f <name of the configuration file (PATH/Configfile)> -c <name of the sops configuration file (PATH/SopsConfigfile)>  
```
create will respond some details 
- static egress ip 
- created managed identity

### Destroy
```
copsctl build destroy -f <name of the configuration file (PATH/Configfile)> -c <name of the sops configuration file (PATH/SopsConfigfile)>
```