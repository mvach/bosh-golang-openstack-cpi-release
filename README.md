# bosh-golang-openstack-cpi-releases

## Build it

```bash
go build -o bin/cpi
```

## Use it

```bash
echo "{\"method\":\"info\",\"arguments\":[],\"context\":{\"director_uuid\":\"ad0d986f-7712-439d-8e07-eb09602239e4\",\"request_id\":\"cpi-269726\",\"vm\":{\"stemcell\":{\"api_version\":3}}}}" | ./bin/cpi -configFile=local/default_config.json
```