output "public_ip" {
  value = tencentcloud_instance.cvm[0].public_ip
}