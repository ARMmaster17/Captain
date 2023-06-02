<template>
  <p>RESPONSE: {{apiVal}}</p>
  <form v-on:submit.prevent="pingAPI">
    <button class="btn btn-primary">PING</button>
  </form>
</template>

<script>
import api from "../services/api";

export default {
  name: "Ping.vue",
  data() {
    return {
      apiVal: '',
    }
  },
  methods: {
    pingAPI() {
        api.post("util/ping", {
              value: "PONG"
          })
          .then((response) => {
            this.apiVal = response.data.value;
          })
          .catch((error) => {
            window.alert(`The API returned an error: ${error}`);
          })
    }
  }
}

</script>

<style scoped>

</style>