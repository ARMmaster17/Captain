class MachineChannel < ApplicationCable::Channel
  private def subscribed
    stream_from "machine_channel"
  end
end