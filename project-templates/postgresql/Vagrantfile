# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.define "ubuntu" do |ubuntu|
    ubuntu.vm.box = "phusion/ubuntu-14.04-amd64"
    ubuntu.vm.provision :shell, path: "scripts/provisioning/vagrant-ubuntu.sh"
    ubuntu.vm.network :private_network, type: :static, ip: "192.168.50.240"
  end

  config.vm.define "centos" do |centos|
    # centos.gui = true
    centos.vm.box = "metcalfc/centos70-docker"
    centos.vm.provision :shell, path: "scripts/provisioning/vagrant-centos.sh"
    centos.vm.network :private_network, type: :static, ip: "192.168.50.241"
  end

  config.ssh.insert_key = 'true'

  config.ssh.forward_agent = true

  config.vm.synced_folder ENV['GOPATH'], "/go"

  # To connect from vagrant guest OS to your host OS, use the gateway IP address.
  #
  # To find your gateway IP address, run this command:
  # netstat -rn | grep "^0.0.0.0 " | cut -d " " -f10
  #
  # Here's an example on how to run your application inside vagrant while connecting to your host's PostgreSQL.
  # GOPATH=/go DSN=postgres://$PGUSER@$(netstat -rn | grep "^0.0.0.0 " | cut -d " " -f10):5432/$PROJECT_NAME?sslmode=disable go run main.go
end
