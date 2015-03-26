# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|
  config.vm.provider "virtualbox"
  config.vm.provider "vmware_fusion"

  config.env.enable

  config.vm.box = "debian720"
  config.vm.box_check_update = false

  config.vm.network "private_network", ip: "192.168.33.31"

  config.vm.hostname = "khabar"

  config.vm.provider "virtualbox" do |v|
    v.memory = 256
    v.cpus = 1
  end

  enable_puppet_script = <<-EOF
    puppet agent --enable
  EOF

  config.vm.provision "shell" do |s|
    s.inline = enable_puppet_script
  end

end
