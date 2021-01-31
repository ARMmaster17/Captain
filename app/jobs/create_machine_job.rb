require 'erb'

include LxcLib
include ERB::Util

class CreateMachineJob < ApplicationJob
  queue_as :default

  def perform(vm)
    #MachineChannel.broadcast_to(vm, "#{vm.hostname} is being created...")
    ActionCable.server.broadcast("machine_channel", body: "#{vm.hostname} is being created...")
    mod_ip_address = vm.ip_address.gsub("/", "%2F")
    vmid = LxcLib.get_vmid
    vm_settings = {}
    vm_settings['vmid'] = vmid
    vm_settings['hostname'] = vm.hostname
    vm_settings['cores'] = vm.cpu
    vm_settings['memory'] = vm.ram
    vm_settings['swap'] = vm.ram
    #vm_settings['unique'] = '1'
    vm_settings['net0'] = "name%3Deth0%2Cbridge%3Dinternal%2Cip%3D#{mod_ip_address}%2Cgw%3D10.0.0.1%2Cfirewall%3D0%2Cmtu%3D1450"
    vm_settings['nameserver'] = "10.0.0.1"
    vm_settings['storage'] = 'pve-storage'
    vm_settings['rootfs'] = "pve-storage%3A#{vm.disk}"
    vm_settings['ostemplate'] = 'pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz'
    vm_settings['onboot'] = '1'
    #vm_settings['ssh-public-keys'] = ''
    vm_settings['start'] = '1'
    vm_settings['unprivileged'] = '1'
    packed_settings = vm_settings.to_a.map { |v| v.join '=' }.join '&'

    LxcLib.create_machine(packed_settings)

    vm.vmid = vmid
    vm.save

    #MachineChannel.broadcast_to(vm, "#{vm.hostname} is now online.")
    ActionCable.server.broadcast("machine_channel", body: "#{vm.hostname} is now online.")
  end
end
