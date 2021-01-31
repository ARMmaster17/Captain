import consumer from "./consumer"

consumer.subscriptions.create({ channel: "MachineChannel"}, {
    connected() {
        console.log("ActiveCable is working!");
    },
    received(data) {
        console.log(data["body"]);
    }
})