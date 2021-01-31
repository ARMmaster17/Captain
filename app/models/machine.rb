include LxcLib

class Machine < ApplicationRecord
  after_commit :create_machine, on: :create
  before_destroy :destroy_machine

  private
  def create_machine
    CreateMachineJob.perform_later(self)
    ActionCable.server.broadcast("machine_channel", body: "#{self.hostname} is queued for provisioning.")
  end

  def destroy_machine
    ActionCable.server.broadcast("machine_channel", body: "#{self.hostname} is queued for deletion.")
    LxcLib.delete_machine(self.vmid)
  end
end
