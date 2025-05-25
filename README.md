# mzmetrics
Prometheus scrape endpoint for MyZone airconditioners

## Run with docker compose

### MZMetrics
Create a .env file with the endpoint for your AC controller. Remplace <controller_host> with the 
address of your AC controller.
```
echo MZMETRICS_AC_URL=http://<controller_host>:2025/getSystemData > .env
```

### Prometheus
docker-compose is setup to read a `prometheus.yml` from the repo root.
```
cp prometheus.example.yml prometheus.yml`
```
The configure remote_write for your remote endpoint. Docker-compose is configured to load the 
password from a mounted secret file. Write your secret to the `remote-write-password` file.

Example using OSX keychain:
```
security find-generic-password -s <secret_name> -w > remote-write-password
```

Start the stack with docker-compose
```
docker-compose up
```
