# Create service account for staging
resource "google_service_account" "tuneslap_api" {
  account_id   = "tuneslap-api-staging"
  display_name = "Tuneslap Staging API Service Account"
  description  = "Service account for Tuneslap API in staging environment."

  lifecycle {
    prevent_destroy = true
  }
}

# Logs bucket for access logging
resource "google_storage_bucket" "logs_bucket" {
  name          = "tuneslap-logs-staging"
  location      = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket" "media_bucket" {
  name          = "tuneslap-media-staging"
  location      = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  # Media bucket needs public access for serving processed audio/images to users
  # Files are served via direct URLs for low-latency playback

  versioning {
    enabled = true
  }

  logging {
    log_bucket = google_storage_bucket.logs_bucket.name
  }

  cors {
    origin          = ["http://localhost:3001", "https://*.tuneslap.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "OPTIONS"]
    response_header = ["Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Max-Age"]
    max_age_seconds = 3600
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Media bucket requires public read access for serving processed audio/images
# This is intentional - media files need to be publicly accessible for playback
# tfsec:ignore:google-storage-no-public-access
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.media_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_storage_bucket" "system_bucket" {
  name          = "tuneslap-system-staging"
  location      = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }

  logging {
    log_bucket = google_storage_bucket.logs_bucket.name
  }

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket" "user_uploads_bucket" {
  name          = "tuneslap-user-uploads-staging"
  location      = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  versioning {
    enabled = true
  }

  logging {
    log_bucket = google_storage_bucket.logs_bucket.name
  }

  cors {
    origin          = ["http://localhost:3001", "https://*.tuneslap.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "OPTIONS"]
    response_header = ["Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Max-Age"]
    max_age_seconds = 3600
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Allow API service account to write to both buckets
resource "google_storage_bucket_iam_member" "api_write" {
  for_each = toset([
    google_storage_bucket.media_bucket.name,
    google_storage_bucket.system_bucket.name,
    google_storage_bucket.user_uploads_bucket.name
  ])
  bucket = each.key
  role   = "roles/storage.objectCreator"
  member = "serviceAccount:${google_service_account.tuneslap_api.email}"
}

# Allow API service account to read from both buckets
resource "google_storage_bucket_iam_member" "api_read" {
  for_each = toset([
    google_storage_bucket.media_bucket.name,
    google_storage_bucket.system_bucket.name,
    google_storage_bucket.user_uploads_bucket.name
  ])
  bucket = each.key
  role   = "roles/storage.objectViewer"
  member = "serviceAccount:${google_service_account.tuneslap_api.email}"
} 