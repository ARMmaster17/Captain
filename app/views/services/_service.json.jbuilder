json.extract! service, :id, :name, :scale, :cpu, :ram, :disk, :hostname, :domain, :created_at, :updated_at
json.url service_url(service, format: :json)
