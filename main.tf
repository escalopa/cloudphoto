
terraform {
  required_providers {
    yandex = {
      source = "yandex-cloud/yandex"
    }
  }
}

provider "yandex" {
  token     = var.access_token
  cloud_id  = var.cloud_id
  folder_id = var.folder_id
  zone      = "ru-central1-a"
}

resource "yandex_iam_service_account" "sa" {
  name        = "${var.app_id}-bucket-manager"
  description = "bucket manager service account"
  folder_id   = var.folder_id
}

resource "yandex_resourcemanager_folder_iam_member" "sa-storage-admin" {
  folder_id = var.folder_id
  role      = "storage.admin"
  member    = "serviceAccount:${yandex_iam_service_account.sa.id}"
}

resource "yandex_iam_service_account_static_access_key" "sa-static-key" {
  service_account_id = yandex_iam_service_account.sa.id
  description        = "static access key for object storage"
}

resource "yandex_storage_bucket" "app-storage" {
  access_key = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  secret_key = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  bucket     = "${var.app_id}-storage"
  default_storage_class = "STANDARD"
  max_size   = 1073741824
}

output "sa-static-key-access-key" {
  value = yandex_iam_service_account_static_access_key.sa-static-key.access_key
  sensitive = true
}

output "sa-static-key-secret-key" {
  value = yandex_iam_service_account_static_access_key.sa-static-key.secret_key
  sensitive = true
}