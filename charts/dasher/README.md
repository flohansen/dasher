# Dasher Helm Chart

## Install

```bash
helm install <release-name> 
```

## Configuration
| Parameter | Description | Default |
| --------- | ----------- | ------- |
| replicas  | It is not recommended to set this value different to 1, because it is not possible to change the database type for now. Setting the value to n > 1 would result in n pods with each having its own SQLite database attached.  Replication using SQLite is not supported yet. Anyways, this should not be an issue since the load on the application will not affect the server. | `1` |
| ingress.enabled | If an ingress should be deployed. Make sure to setup a proper ingress controller e.g. [nginx](https://docs.nginx.com/nginx-ingress-controller/). | `false` |
| ingress.ingressClassName | The ingress class name (like `nginx`) | `""` |
| storage.size | The size of the persistent volume. This is being used to save the application state (in SQLite). | `1G` |