require 'erb'

include LxcLib
include ERB::Util

class CreateMachineJob < ApplicationJob
  queue_as :default

  def perform(vm)
    vmid = LxcLib.get_vmid
    vm_settings = {}
    vm_settings['vmid'] = vmid
    vm_settings['hostname'] = vm.hostname
    vm_settings['cores'] = vm.cpu
    vm_settings['memory'] = vm.ram
    vm_settings['swap'] = vm.ram
    #vm_settings['unique'] = '1'
    vm_settings['net0'] = url_encode("name=eth0,bridge=internal,ip=#{vm.ip_address},gw=10.0.0.1,ip6=manual,firewall=0,mtu=1450")
    vm_settings['storage'] = 'pve-storage'
    vm_settings['rootfs'] = "0,size=#{vm.disk}"
    vm_settings['ostemplate'] = 'pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz'
    vm_settings['onboot'] = '1'
    #vm_settings['ssh-public-keys'] = ''
    vm_settings['start'] = '1'
    vm_settings['unprivileged'] = '1'
    packed_settings = vm_settings.to_a.map { |v| v.join '=' }.join '&'

    Rails.logger.warn packed_settings.to_s

    LxcLib.create_machine(packed_settings)

    vm.vmid = vmid
    vm.save
  end
end
