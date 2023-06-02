<template>
<ul>
    <li v-for="flight in flights">
        <router-link :to="`/rsx/flight/${flight.ID}`" class="nav-link">{{ flight.name }}</router-link>
    </li>
</ul>
</template>

<script>
import FlightService from "../services/flight.service";

export default {
    name: "Flights",
    data() {
        return {
            flights: [],
        }
    },
    mounted() {
        FlightService.getFlights().then(
            (response) => {
                this.flights = response.data;
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