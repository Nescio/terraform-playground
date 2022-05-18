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
  name    = lower( join("", regexall("[a-zA-Z]+", "1Cha 2rl_i1e")))
  species = "dog"
  age     = 7
}
