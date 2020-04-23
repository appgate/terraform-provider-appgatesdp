terraform {
  required_version = ">= 0.12"

  required_providers {
    aws = ">= 2.45.0"
  }
}

provider "aws" {
  region = var.aws_region
}


resource "aws_subnet" "appgate_appliance_subnet" {
  vpc_id     = var.vpc_id
  cidr_block = var.appliance_cidr_block

  assign_ipv6_address_on_creation = false

  tags = {
    Name   = "Appgate-Subnet"
    Group  = "Appgate"
    Author = "dln"
  }
}

resource "aws_route_table_association" "appgate_route_table_assoication" {
  subnet_id      = aws_subnet.appgate_appliance_subnet.id
  route_table_id = aws_route_table.appgate_route_table.id
}

data "aws_internet_gateway" "selected" {
  internet_gateway_id = var.internet_gateway_id
}

resource "aws_route_table" "appgate_route_table" {
  vpc_id = var.vpc_id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = data.aws_internet_gateway.selected.id
  }



  tags = {
    Name   = "AppgateRouting"
    Group  = "Appgate"
    Author = "dln"
  }
}

resource "aws_security_group" "appgate_security_group" {
  description = "Security group used by AppGate."
  vpc_id      = var.vpc_id

  # Allow all protocols
  ingress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"

    cidr_blocks = [
      var.appliance_cidr_block,
      "212.16.176.132/32",
      "2.248.12.29/32"
    ]
    ipv6_cidr_blocks = [
      "2a01:2b0:302c::/48"
    ]

  }

  ingress {
    from_port = 443
    to_port   = 443
    protocol  = "tcp"

    cidr_blocks = [
      var.appliance_cidr_block,
    ]
  }

  egress {
    from_port = 0
    to_port   = 0
    protocol  = "-1"

    cidr_blocks = [
      "0.0.0.0/0",
    ]
  }

  tags = {
    Name   = "Appgate-Security-Group"
    Group  = "Appgate"
    Author = "dln"
  }
}

data "aws_ami" "appgate" {
  most_recent = true

  filter {
    name   = "name"
    values = ["AppGate-SDP-5.1.*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["194335033129"] # Canonical
}

resource "aws_key_pair" "deployer" {
  key_name   = "deployer-key"
  public_key = file(var.public_key)

  tags = {
    Name   = "controller-appgate"
    Group  = "Appgate"
    Author = "dln"
  }
}


resource "aws_instance" "appgate_controller" {
  ami = data.aws_ami.appgate.id
  # https://sdphelp.cyxtera.com/adminguide/v5.0/instance-sizing.html
  instance_type = var.controller_instance_type

  subnet_id = aws_subnet.appgate_appliance_subnet.id
  vpc_security_group_ids = [
    aws_security_group.appgate_security_group.id
  ]

  key_name = aws_key_pair.deployer.key_name


  associate_public_ip_address = true


  root_block_device {
    volume_type = "gp2"
    volume_size = 50
  }


  ebs_block_device {
    volume_type = "gp2"
    volume_size = 20
    device_name = "/dev/xvdb"
  }


  connection {
    type        = "ssh"
    user        = "cz"
    timeout     = "25m"
    private_key = file(var.private_key)
    host        = aws_instance.appgate_controller.public_ip
  }



  # https://sdphelp.cyxtera.com/adminguide/v5.0/appliance-installation.html
  provisioner "remote-exec" {
    inline = [
      # https://sdphelp.cyxtera.com/adminguide/v5.0/new-appliance.html?anchor=manual-seeding
      "cz-seed --output /home/cz/seed.json --password cz cz --dhcp-ipv4 eth0 --enable-logserver --no-registration --hostname ${self.public_dns}"
    ]
  }
  user_data_base64 = "IyEvdXNyL2Jpbi9lbnYgYmFzaAoKdG91Y2ggL3RtcC9zZWVkZWQgJj4gL2Rldi9udWxsCg=="


  tags = {
    Name   = "controller-appgate"
    Group  = "Appgate"
    Author = "dln"
  }
}

resource "aws_instance" "appgate_gateway" {
  ami = data.aws_ami.appgate.id
  # https://sdphelp.cyxtera.com/adminguide/v5.0/instance-sizing.html
  instance_type = var.controller_instance_type

  subnet_id = aws_subnet.appgate_appliance_subnet.id
  vpc_security_group_ids = [
    aws_security_group.appgate_security_group.id
  ]

  key_name = aws_key_pair.deployer.key_name


  associate_public_ip_address = true


  root_block_device {
    volume_type = "gp2"
    volume_size = 50
  }


  ebs_block_device {
    volume_type = "gp2"
    volume_size = 20
    device_name = "/dev/xvdb"
  }


  user_data_base64 = base64encode("cz-seed --output /home/cz/seed.json --password cz cz --dhcp-ipv4 eth0 --no-registration --disable-controller")


  tags = {
    Name   = "gateway-appgate"
    Group  = "Appgate"
    Author = "dln"
  }
}
