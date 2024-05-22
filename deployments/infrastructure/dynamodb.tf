resource "aws_dynamodb_table" "products" {
  name         = "nr-memory-leak-investigation"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "ProductKey"

  attribute {
    name = "ProductKey"
    type = "S"
  }
}
