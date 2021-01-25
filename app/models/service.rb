class Service < ApplicationRecord
  has_many :machines

  after_commit :create_machines, on: :create
  before_commit :destroy_machines, on: :destroy
end

def create_machines
  (Range.new(1, self.scale)).each do |n|
    newmachine = Machine.create(
      hostname: "#{self.hostname}#{n}.#{self.domain}",
      ip_address: "10.0.0.0/16",
      vmid: -1,
      cpu: self.cpu,
      ram: self.ram,
      disk: self.disk,
      service_id: self.id
    )
    sleep(1)
  end
end

def destroy_machines
  Machine.where(service_id: self.id) do |m|
    m.destroy
  end
end
