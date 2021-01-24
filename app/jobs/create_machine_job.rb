include LxcLib

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
    #vm_settings['net0'] = '"name=eth0,bridge=internal,ip=dhcp"'
    vm_settings['storage'] = 'pve-storage'
    vm_settings['ostemplate'] = 'pve-img:vztmpl/debian-10-standard_10.7-1_amd64.tar.gz'
    packed_settings = vm_settings.to_a.map { |v| v.join '=' }.join '&'

    LxcLib.create_machine(packed_settings)

    vm.vmid = vmid
    vm.save
    while(true)
      begin
        LxcLib.start_machine(vmid)
        return
      rescue
        sleep(5)
      end
    end

  end
end
