# Kate
## Summary
Turns out if you throw CoreOS Clair into your Kubernetes namespace with the help of a friend Kate will automatically scan all new launched containers. Kate will also rescan all the images every couple of hours just to let you know if the CVE situation has changed.

## Example Deployment
- Start an example NGINX container
```bash
kubectl run nginx --image=nginx --restart=Always
```
- Deploy the ./example/cve-storage.yaml PVC to persistance for your Clair CVE tracker
```bash
kubectl create -f ./example/cve-storage.yaml
```
- Deploy the Clair configuration and Postgres database
```bash
kubectl create secret generic clairsecret --from-file=./example/config.yaml
kubectl create -f ./example/clair-postgres.yaml
```
- Deploy Clair and Kate
```bash
kubectl create -f ./example/clair-postgres.yaml
```

- Find the pod name for Clar and Kate
```bash
watch kubectl get pods -o wide
```

- Monitor and wait for the Clair CVE crawl to finish. Look for something like an 'updater: update finished' message.
```bash
kubectl logs clair-and-kate-66ayv -c clair -f
```

- Access your kate service on http://127.0.0.1:8080
```bash
kubectl port-forward clair-and-kate-1jzhx 8080
```

## TODO
Build an additional UI Visualization
