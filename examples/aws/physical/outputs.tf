output "controller_dns" {
  value = aws_instance.appgate_controller.public_dns
}
output "gateway_dns" {
  value = aws_instance.appgate_gateway.public_dns
}
