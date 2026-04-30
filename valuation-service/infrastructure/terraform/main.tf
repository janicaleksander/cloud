terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "6.36.0"
    }
  }
}
provider "aws" {
  region = "us-east-1"
}

resource "aws_db_instance" "main" {
  identifier = "valuation-db"
  engine = "postgres"
  engine_version = "17"
  instance_class = "db.t3.micro"
  allocated_storage = 5

  db_name = "postgres"
  username = "postgres"
  password = var.db_password
  skip_final_snapshot = true
  publicly_accessible = true

  vpc_security_group_ids = [aws_security_group.rds-sg.id]
  tags = {
    Name = "valuation-rds"
  }
}
