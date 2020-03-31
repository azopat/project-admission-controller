# Project access Admission Controller

This project implements a basic admission controller webhook to validate projects.
It uses kubernetes sig package for the implementation

The purpose of this webhook is to limit the amount of project that can be created by a user. This can be usefull on multi-tenant clusters, to limit the amount of project a customer can create.
Therefore, on every project creation request, a validation process will be initiated by the api server and sent to the webhook.
An environment variable will need to be defined in webhook deployment configuration to define the allowed count

## Build

`docker build -t tag .` Should be all that is necessary to build.

## Deploy
PS : The provided deployment script is tested on openshift 4.x cluster, you may need to adjust it a bit to have it running on other kubernetes cluster

- Build the docker image as described above and push it to a registry accessible by your cluster
- Edit the deploy.sh script setting the variables
- Run the script, which will create and configure the ressources needed.
