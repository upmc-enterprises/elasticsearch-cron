# elasticsearch-cron

CronJob container to perform snapshots of an elasticsearch cluster

## Args

The container takes the following args:

- action: Action to perform (repository or snapshot)
- repo-type: Type of repository, s3 or gcs.
- bucket-name: Name of s3 or gcs bucket
- elastic-url: Full dns url to elasticsearch
- auth-username: Authentication username (if applicable)
- auth-password: Authentication password (if applicable)

# About

Built by UPMC Enterprises in Pittsburgh, PA. http://enterprises.upmc.com/