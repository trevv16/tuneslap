# TuneSlap Infrastructure

This directory contains the Terraform configurations for TuneSlap infrastructure. We maintain separate configurations for staging and production environments in different Google Cloud projects.

## Project Structure

```
infra/
├── README.md
├── production/
│   ├── main.tf
│   └── providers.tf
└── staging/
    ├── main.tf
    └── providers.tf
```

## Prerequisites

1. Install [Terraform](https://www.terraform.io/downloads.html) (v1.0.0 or later)
2. Install [Google Cloud SDK](https://cloud.google.com/sdk/docs/install)
3. Configure Google Cloud authentication:
   ```bash
   gcloud auth application-default login
   ```

## Google Cloud Projects

- Production: `munene-dev`
- Staging: `munene-dev-staging`

## Commands

```bash
# Navigate to staging/production directory
cd infra/staging
or
cd infra/production

# Initialize Terraform
terraform init

# Plan the changes
terraform plan

# Apply the changes
terraform apply

# Destroy resources (if needed)
terraform destroy
```

## Resources Created

### Production

- Service Account: `tuneslap-api@munene-dev.iam.gserviceaccount.com`
- System Bucket: `tuneslap-system-prod`
- Media Bucket: `tuneslap-media-prod`

### Staging

- Service Account: `tuneslap-api-staging@munene-dev-staging.iam.gserviceaccount.com`
- System Bucket: `tuneslap-system-staging`
- Storage Bucket: `tuneslap-media-staging`

## Bucket Access Policies

### Public Buckets (Intentional)

The following buckets have public read access (`allUsers` with `roles/storage.objectViewer`):

| Bucket | Purpose | Why Public |
|--------|---------|------------|
| `tuneslap-media-{env}` | Processed audio/images | Direct playback in browser without signed URLs |
| `tuneslap-user-uploads-{env}` | User-uploaded media | Direct playback in browser without signed URLs |

**Controls in place:**

- Access logging enabled (logs written to `tuneslap-logs-{env}`)
- CORS restricted to `tuneslap.com` and subdomains (production) or localhost (staging)
- Uniform bucket-level access enforced (no per-object ACLs)
- All objects are user-generated soundboard content only

### Private Buckets

| Bucket | Purpose | Access |
|--------|---------|--------|
| `tuneslap-system-{env}` | Internal system data | Service account only |
| `tuneslap-logs-{env}` | Access logs | Service account only, public access prevention enforced |

### Service Account Permissions

The API service account has:

- `roles/storage.objectCreator` on media, system, and user uploads buckets
- `roles/storage.objectViewer` on media, system, and user uploads buckets

## Best Practices

1. Always run `terraform plan` before `apply`
2. Review the plan output carefully
3. Use separate service accounts for each environment
4. Keep production and staging configurations separate
5. Never modify cloud resources directly

## Troubleshooting

If you encounter authentication issues:
```bash
# Re-authenticate with Google Cloud
gcloud auth application-default login

# Verify your current project
gcloud config get-value project
```

## Security Notes

- Production credentials should be kept secure
- Never commit service account keys to version control
- Use environment variables for sensitive values
- Regularly rotate service account keys 