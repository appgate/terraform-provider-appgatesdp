variable "aws_region" {
  description = "AWS region to launch servers."
  default     = "us-east-1"
  # default     = "eu-north-1"
}

variable "vpc_id" {
  description = "VPC used for appgate."
  default     = "vpc-1e9b5879" # envy vpc
  # default     = "vpc-1a5bb673" # eu-north-1
}

variable "private_key" {
  description = "Path to SSH private key that is used with key_name."
}
variable "public_key" {
  description = "Path to the public key"
}


variable "appliance_cidr_block" {
  description = "The network addresses used for appliances."
  default     = "10.0.242.0/24"
  # default     = "172.31.252.0/24" # stockholm
}
variable "internet_gateway_id" {
  default = "igw-e2dfaa86"
}

variable "controller_instance_type" {
  description = "aws instance size for the Controller. See https://sdphelp.cyxtera.com/adminguide/v5.0/instance-sizing.html"
  default     = "m4.xlarge"
}

variable "gateway_instance_type" {
  description = "aws instance size for the Controller. See https://sdphelp.cyxtera.com/adminguide/v5.0/instance-sizing.html"
  default     = "m4.xlarge"
}
