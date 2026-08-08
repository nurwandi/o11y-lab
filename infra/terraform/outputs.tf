output "public_ip" {
  description = "Public IP of the k3s node"
  value       = aws_instance.node.public_ip
}

output "ssh_command" {
  description = "SSH into the node"
  value       = "ssh ubuntu@${aws_instance.node.public_ip}"
}

output "kubeconfig_hint" {
  description = "How to get a local kubeconfig"
  value       = "scp ubuntu@${aws_instance.node.public_ip}:/etc/rancher/k3s/k3s.yaml ./kubeconfig && sed -i '' 's/127.0.0.1/${aws_instance.node.public_ip}/' ./kubeconfig"
}

output "grafana_hint" {
  description = "Reach Grafana via an SSH tunnel, then open http://localhost:3001"
  value       = "kubectl -n platform port-forward svc/grafana 3001:3000  # (with kubeconfig above)"
}
