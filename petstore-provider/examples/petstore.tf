terraform {
  required_providers {
    petstore = {
      version = "~> 1.0.0"
      source  = "hacknights.club/petstore-provider/petstore"
    }
  }
}

provider "petstore" {
  address = "http://localhost:8080"
}

resource "petstore_pet" "my_pet" {
  name    = "snowball"
  species = "cat"
  age     = 7
}
