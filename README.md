# Deploying the API Gateway on GCP

Follow the instructions [here](https://cloud.google.com/api-gateway/docs/get-started-cloud-run-grpc#creating_an_api_config_with_grpc).

```bash
source .venv/bin/activate
```

```python
python3 -m grpc_tools.protoc \
    --include_imports \
    --include_source_info \
    --proto_path=./proto \
    --descriptor_set_out=./proto/api_descriptor.pb \
    --go_out=. \
    --go-grpc_out=. \
    api.proto
```

Set the project.

```bash
gcloud config set project qiba-core-441819
```

Create/update the API Gateway config. Increment the semvar (0-1-0). E.g: 0-2-0.

```bash
gcloud api-gateway api-configs create grpc-config-0-1-0 \
--api=qiba --project=qiba-core-441819 \
--grpc-files=./proto/api_descriptor.pb,./proto/api_config.yaml
```

Create/update the API Gateway. Match the semvar (0-1-0). E.g: 0-2-0.

Creating...

```bash
gcloud api-gateway gateways create qiba \
  --api=qiba --api-config=grpc-config-0-1-0 \
  --location=us-central1 --project=qiba-core-441819
```

Updating...

```bash
gcloud api-gateway gateways update qiba \
  --api=qiba --api-config=grpc-config-0-2-0 \
  --location=us-central1 --project=qiba-core-441819
```

Describe the service.

```bash
gcloud run services describe qiba
```

Enable HTTP/2.

```bash
gcloud run services update qiba --use-http2
```

Deploy source code.

```bash
gcloud run deploy --source .
```

# Depoying the GAME on GCP

Follow the instructions [here](https://cloud.google.com/run/docs/quickstarts/frameworks/deploy-nextjs-service).
