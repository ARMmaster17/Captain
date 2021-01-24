include LxcLib

class Machine < ApplicationRecord
  after_commit :create_machine, on: :create
  before_destroy :destroy_machine

  private
  def create_machine
    CreateMachineJob.perform_now(self)
  end

  def destroy_machine
    LxcLib.delete_machine(self.vmid)
  end
end
