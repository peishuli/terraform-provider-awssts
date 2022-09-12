terraform {
  required_providers {
    awssts = {
      version = "0.1"
      source  = "peishuli.com/dev/awssts"
    }
  }
}

provider "awssts" {
  user_name = "sam" 
}

data "aws_federation_token" "current" { 
  provider = awssts
}

output user_id {
  value = data.aws_federation_token.current.federated_user_id
}

output federation_token {
  value = jsondecode(data.aws_federation_token.current.federation_token)
}

