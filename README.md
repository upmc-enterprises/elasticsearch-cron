# elasticsearch-cron

CronJob container to perform snapshots of an elasticsearch cluster

## Args

The container takes the following args:

- action: Action to perform (repository or snapshot)
- s3-bucket-name: Name of s3 bucket
- elastic-url: Full dns url to elasticsearch
- auth-username: Authentication username (if applicable)
- auth-password: Authentication password (if applicable)

# About

Built by UPMC Enterprises in Pittsburgh, PA. http://enterprises.upmc.com/