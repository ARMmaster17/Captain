import consumer from "./consumer"
import toastr from "toastr"

consumer.subscriptions.create({ channel: "MachineChannel"}, {
    connected() {
        console.log("ActiveCable is working!");
        toastr.info("ActiveCable is working!");
    },
    received(data) {
        console.log(data["body"]);
        toastr.info(data["body"]);
    }
})