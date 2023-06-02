import { createWebHistory, createRouter } from "vue-router";
import Login from "./components/Login.vue";
import Home from "./components/Home.vue";
import Ping from "./components/Ping.vue";
import Deployments from "./components/Deployments.vue";
import Flights from "./components/Flights.vue";
import Flight from "./components/Flight.vue";

const routes = [
    {
        path: "/",
        name: "home",
        component: Home,
    },
    {
        path: "/login",
        component: Login,
    },
    {
        path: "/ping",
        component: Ping
    },
    {
        path: "/deployments",
        component: Deployments
    },
    {
        path: "/rsx/flights",
        component: Flights
    },
    {
        path: "/rsx/flight/:id",
        component: Flight
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;