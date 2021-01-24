require 'rest-client'
require 'json'

module LxcLib
  def LxcLib.create_machine(vm_settings)
    resource = get_resource
    auth_header = authenticate(resource)

    begin
      return resource['nodes/pxvh1/lxc'].post(vm_settings, auth_header)
    rescue RestClient::ExceptionWithResponse => e
      Rails.logger.warn e.http_body
    end
  end

  def LxcLib.delete_machine(vmid)
    resource = get_resource
    auth_header = authenticate(resource)

    stop_machine(vmid)
    sleep(3)
    begin
      resource['nodes/pxvh1/lxc/' + String(vmid)].delete(auth_header.merge({'force' => '1', 'purge' => '1'}))
    rescue
      Rails.logger.warn "Unable to delete #{vmid}, ignoring..."
    end
  end

  def LxcLib.start_machine(vmid)
    resource = get_resource
    auth_header = authenticate(resource)

    wait_for_machine_ready(vmid)

    resource["nodes/pxvh1/lxc/#{vmid}/status/start"].post({}, auth_header)
  end

  def LxcLib.stop_machine(vmid)
    resource = get_resource
    auth_header = authenticate(resource)

    resource["nodes/pxvh1/lxc/#{vmid}/status/stop"].post({}, auth_header)
  end

  def LxcLib.wait_for_machine_ready(vmid)
    resource = get_resource
    auth_header = authenticate(resource)

    while(true)
      sleep(5)
      data = resource["nodes/pxvh1/lxc/#{vmid}/status/current"].get(auth_header).body
      locks = JSON.parse(data)['data']['lock']
      Rails.logger.warn "Value of locks: '#{locks}'"
      if(locks.eql?(""))
        sleep(5)
        return
      end
    end
  end

  def LxcLib.get_vmid
    resource = get_resource
    auth_header = authenticate(resource)

    vmid = resource['cluster/nextid'].get(auth_header).body
    Rails.logger.warn('Got VMID: ' + vmid)
    data = JSON.parse(vmid)
    return data['data']
  end

  def LxcLib.get_resource
    return RestClient::Resource.new('https://192.168.1.241:8006/api2/json/', {verify_ssl: false})
  end
  def LxcLib.authenticate(resource)
    payload = { username: ENV['PROXMOX_USER'], password: ENV['PROXMOX_PASSWORD'] }
    response = resource['access/ticket'].post(payload)
    data = JSON.parse(response.body)
    ticket = data['data']['ticket']
    csrf = data['data']['CSRFPreventionToken']
    token = 'PVEAuthCookie=' + + ticket.gsub!(/:/, '%3A').gsub!(/=/, '%3D')
    return { CSRFPreventionToken: csrf, cookie: token, 'ContentType' => 'application/x-www-form-urlencoded' }
  end
end