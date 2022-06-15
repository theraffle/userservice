# Deployment Guide

This is a guide to deploy UserService on your cluster. The contents are as follows.

* [Prerequisites](#prerequisites)
* [Deploy UserService](#deploy-userservice)

## Prerequisites
- [MySQL Database](./tables.md)
- Kubernetes Cluster

## Deploy UserService
1. Create secret for MySQL database info (Use `secret.yaml` in samples directory)
   ```yaml
   apiVersion: v1
   kind: Secret
   metadata:
     name: db-secret
     namespace: default
   stringData:
   # Fill in your DB info
     dbHost: ""
     dbPort: ""
     dbUser: ""
     dbPassword: ""
     dbName: ""
   type: Opaque
   ```
   ```bash
   kubectl apply -f ./samples/secret.yaml
   ```
   
2. Apply `userservice.yaml`
   ```bash
   kubectl apply -f ./kubernetes-manifests/release.yaml
   ```
3. Wait until `userservice` pod is ready
   ```bash
   kubectl get pod -w
   ```