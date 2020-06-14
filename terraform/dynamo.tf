resource "aws_dynamodb_table" "oops_table" {
  name           = "oops"
  billing_mode   = "PROVISIONED"
  read_capacity  = 5
  write_capacity = 5
  hash_key       = "OopsID" // Do not change

  attribute {
    name = "OopsID"
    type = "S"
  }

  ttl {
    attribute_name = "Expiration" // Do not change
    enabled        = true
  }


  tags = {
    Name        = "oops"
    Environment = "prod"
    ManagingTeam = "SysEng"
    CostCenter = "SysEng"
  }
}
