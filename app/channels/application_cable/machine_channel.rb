class MachineChannel < ApplicationCable::Channel
  def subscribed
    stream_from "machine_stream"
  end
end