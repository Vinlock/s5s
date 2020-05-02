# Secrets (s5s)

s5s is a tool to download and apply secrets from cloud Secret Managers

## Supported Secret Managers
- [Google Cloud Secret Manager](https://github.com/Vinlock/s5s#google-cloud-secret-manager) [(link)](https://cloud.google.com/secret-manager)

### Google Cloud Secret Manager
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