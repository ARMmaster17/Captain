json.extract! machine, :id, :hostname, :ip_address, :vmid, :cpu, :ram, :disk, :created_at, :updated_at
json.url machine_url(machine, format: :json)
