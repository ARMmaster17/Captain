<template>
  <h1>Flight Details</h1>
  <p>ID: {{flight.ID }}</p>
  <p>Name: {{ flight.name }}</p>
</template>

<script>
import FlightService from "../services/flight.service";

export default {
    name: "Flight",
    data() {
        return {
            flight: {},
        }
    },
    mounted() {
        FlightService.getFlight(this.$route.params.id).then(
            (response) => {
                this.flight = response.data;
            },
            (error) => {
                console.log(error);
            }
        )
    }
}
</script>

<style scoped>

</style>