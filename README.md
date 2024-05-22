# Memory Leak Investigation

## Infrastructure Setup

### Database

The application reads data from a DynamoDB database. It's the only app that we own that does so, and it's the only one that seems to experience a memory leak. So, we decided to leave this bit as the issue might be in New Relic's instrumentation for the AWS SDK.

We include a terraform script to create the database in `deployments/infrastructure`. Adjust the region and table name as desired, and simply run `terraform init` and `terraform apply`. 

We also include a JS script to populate the database with some sample data. Go into `deployments/db-population` and run:

```shell
npm init
node ./populate.js
```

You might need to set the right AWS credentials to conduct the previous steps.

## Running the App

We ran the application on a kubernetes cluster. We have not included any configuration in this repo, as we believe the exact configuration you would need to use will most likely differ substantially. In any case, here are some details that might be useful:

- We use an ALB to expose the application
- Resource limits:
  - Memory: 500Mi
- Resource requests:
  - cpu: 300m
  - memory: 50Mi

### Environment Variables

The application needs the following environment variables to be specified:

- NRML_TABLENAME: DynamoDB table name
- NRML_AWS_REGION: AWS region where the table is
- NRML_NEWRELIC_LICENSE: New Relic license key

## Generating Traffic

See [this README file](./scripts/load-test/README.md) for details on how to generate traffic for the application.

## Profiling

You can enable pprof by setting the environment variable `NRML_ENABLE_PPROF` to `true`. This uses `gin-contrib/pprof`. See more details [here](https://github.com/gin-contrib/pprof).