# Single EC2 instance running k3s. cloud-init installs k3s + ArgoCD and points
# ArgoCD at this repo's gitops/ directory (GitOps takes over from there).

data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
}

resource "aws_instance" "node" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  subnet_id              = aws_subnet.public.id
  vpc_security_group_ids = [aws_security_group.node.id]
  key_name               = var.key_name

  user_data = templatefile("${path.module}/user_data.sh.tftpl", {
    gitops_repo_url = var.gitops_repo_url
  })

  root_block_device {
    volume_type = "gp3"
    volume_size = var.root_volume_gb
  }

  tags = { Name = "o11y-lab" }
}
