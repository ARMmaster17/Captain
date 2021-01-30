import consumer from "./consumer"

consumer.subscriptions.create({ channel: "MachineChannel", room: "machine_stream"}, {
    received(data) {
        new Notification(data["title"], body: data["body"])
    }
})