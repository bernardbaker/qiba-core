# Deploying the API Gateway on GCP

Follow the instructions [here](https://cloud.google.com/api-gateway/docs/get-started-cloud-run-grpc#creating_an_api_config_with_grpc).

Log into the qiba.fun Google Account. After successful login you should see a page similar to [this](https://cloud.google.com/sdk/auth_success).

```bash
gcloud auth login
```

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
gcloud config set project qiba-core-1-0-0
```

Deploy source code. Service name (press enter). Region (32)

```bash
gcloud run deploy --source . \
  --max-instances=5 \
  --cpu=1 \
  --memory=512Mi
```

Describe the service.

```bash
gcloud run services describe qiba
```

Enable HTTP/2.

```bash
gcloud run services update qiba --use-http2
```

Create/update the API Gateway config. Increment the semvar (0-1-0). E.g: 0-2-0.

In the [api_config.yaml] set the `backend` -> `rules` -> `address` to the Cloud Run service `uc` address. E.g. `qiba-frfefehcjq-uc.a.run.app`.

```bash
gcloud api-gateway api-configs create grpc-config-0-1-0 \
--api=qiba --project=qiba-core-1-0-0 \
--grpc-files=./proto/api_descriptor.pb,./proto/api_config.yaml \
--backend-auth-service-account=qiba-core@qiba-core-1-0-0.iam.gserviceaccount.com
```

Create/update the API Gateway. Match the semvar (0-1-0). E.g: 0-2-0.

Creating...

```bash
gcloud api-gateway gateways create qiba \
  --api=qiba --api-config=grpc-config-0-1-0 \
  --location=us-central1 --project=qiba-core-1-0-0
```

Updating...

```bash
gcloud api-gateway gateways update qiba \
  --api=qiba --api-config=grpc-config-0-2-0 \
  --location=us-central1 --project=qiba-core-1-0-0
```

# Depoying the GAME on GCP

Follow the instructions [here](https://cloud.google.com/run/docs/quickstarts/frameworks/deploy-nextjs-service).

Update the `.env.production` -> `QIBA_CORE_API` to the `qiba-core-1-0-0` API Gateway address. E.g. `qiba-<REDACTED>.uc.gateway.dev`.

Copy the `qiba-core` -> `proto/api.proto` to `qiba-game` -> `proto/api.proto`.

Copy the `qiba-core` -> `proto/api.descriptor.pb` to `qiba-game` -> `proto/api.descriptor.pb`.

When deploying use the following commands:

```bash
gcloud run deploy --source . \
  --max-instances=5 \
  --cpu=1 \
  --memory=512Mi
```

- Accept the default service name (press enter).

- Select the region (32).

- Allow unauthenticated invocations.

See [this](https://cloud.google.com/run/docs/authenticating/public) web page if you receive a `Error: Forbidden Your client does not have permission to get URL / from this server.` message when navigating to the deploy Cloud Run URL for the game UI.

You may need to follow the guidance found [here](https://cloud.google.com/blog/topics/developers-practitioners/how-create-public-cloud-run-services-when-domain-restricted-sharing-enforced) on conditional policies when DRS is enabled.

# QiBA Database

The QIBA Core stores data in memory while working the in the development environment. In production it uses MongoDB more information about the cloud based database can be found [here](https://cloud.google.com/mongodb?hl=en&authuser=1).

During local development the database can be referred to as a [repository](./infrastructure).

- [Game repository](./infrastructure/game_repository_db.go)
- [Leaderboard repository](./infrastructure/leaderboard_repository_db.go)
- [Referral repository](./infrastructure/referral_repository_db.go)
- [User repository](./infrastructure/user_repository_db.go)

In production the current database is MongoDB. The repository files are:

- [Game repository](./infrastructure/game_repository_mongo_db.go)
- [Leaderboard repository](./infrastructure/leaderboard_repository_mongo_db.go)
- [Referral repository](./infrastructure/referral_repository_mongo_db.go)
- [User repository](./infrastructure/user_repository_mongo_db.go)
