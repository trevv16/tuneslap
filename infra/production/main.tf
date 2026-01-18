# Create service account for staging
resource "google_service_account" "tuneslap_api" {
  account_id   = "tuneslap-api-production"
  display_name = "Tuneslap Production API Service Account"
  description  = "Service account for Tuneslap API in production environment."

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket" "media_bucket" {
  name = "tuneslap-media-production"  # Production bucket name
  location = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true

  cors {
    origin          = ["https://tuneslap.com", "https://*.tuneslap.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "OPTIONS"]
    response_header = ["Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Methods", "Access-Control-Allow-Headers", "Access-Control-Max-Age"]
    max_age_seconds = 3600
  }

  lifecycle {
    prevent_destroy = true
  }
}

# Allow public read access to the media bucket
resource "google_storage_bucket_iam_member" "public_read" {
  bucket = google_storage_bucket.media_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

# Allow public read access to the user uploads bucket
resource "google_storage_bucket_iam_member" "user_uploads_public_read" {
  bucket = google_storage_bucket.user_uploads_bucket.name
  role   = "roles/storage.objectViewer"
  member = "allUsers"
}

resource "google_storage_bucket" "system_bucket" {
  name = "tuneslap-system-production"  # Production bucket name
  location = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true

  lifecycle {
    prevent_destroy = true
  }
}

resource "google_storage_bucket" "user_uploads_bucket" {
  name = "tuneslap-user-uploads-production"  # Production bucket name
  location = "us-east1"
  storage_class = "STANDARD"

  uniform_bucket_level_access = true

  cors {
    origin          = ["https://tuneslap.com", "https://*.tuneslap.com"]
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