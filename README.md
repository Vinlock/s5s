# Secrets (s5s)

s5s is a tool to download and apply secrets from cloud Secret Managers

## Supported Secret Managers
- [Google Cloud Secrets Manager](https://github.com/Vinlock/s5s#google-cloud-secrets-manager) [(link)](https://cloud.google.com/secret-manager)
- [AWS Secrets Manager](https://github.com/Vinlock/s5s#aws-secrets-manager) [(link)](https://aws.amazon.com/secrets-manager/)

### Google Cloud Secrets Manager
| Flag                    | Description                                                       | Required | Default Value |
|-------------------------|-------------------------------------------------------------------|----------|---------------|
| `--project \| -p`       | GCP Project Name                                                  | X        |               |
| `--key \| -k`           | GCP Key String (must be provided if `--key-file` is not)          | X        |               |
| `--key-file \| -f`      | GCP Key File (JSON) (must be provided if `--key` is not)          | X        |               |
| `--secret \| -s`        | List of secrets formatted as `<k8s secret key>=<gcp secret name>` | X        |               |
| `--output-secret \| -o` | Name of k8s secret                                                | X        |               |
| `--version \| -v`       | GCP Secret Version                                                |          | latest        |

#### Example:
```bash
$ s5s gcp \
    -p gcp-project-id \
    -f secret.json \
    -s "mysqlusername=app-mysql-username" \
    -s "mysqlpassword=app-mysql-password" \
    -o mysql-creds | kubectl apply --context k8s-cluster -n app-namespace -f -
```

### AWS Secrets Manager
| Flag                    | Description                                                       | Required | Default Value         |
|-------------------------|-------------------------------------------------------------------|----------|-----------------------|
| `--secret \| -s`        | AWS Secret Name                                                   | X        |                       |
| `--region \| -r`        | AWS Region Name                                                   |          | us-west-2             |
| `--accessKeyId`         | AWS Access Key ID                                                 |          | AWS_ACCESS_KEY_ID env |
| `--accessKey`           | AWS Access Key                                                    |          | AWS_ACCESS_KEY env    |
| `--token`               | AWS Access Token                                                  |          |                       |
| `--output-secret \| -o` | Name of k8s secret                                                | X        |                       |

#### Example:
```bash
$ s5s aws \
    -s "project/mysql/secrets"
    --accessKeyId AW12312312412
    --accessKey XpijOIPUYh087^*&(^%
    -o mysql-creds | kubectl apply --context k8s-cluster -n app-namespace -f -
```

### Download Latest
- Linux (https://github.com/Vinlock/s5s/releases/download/linux-latest/s5s)
- Windows (https://github.com/Vinlock/s5s/releases/download/windows-latest/s5s.exe)
- MacOS (https://github.com/Vinlock/s5s/releases/download/darwin-latest/s5s)
