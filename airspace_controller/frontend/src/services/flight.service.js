import api from "./api";

class FlightService {
    getFlights() {
        return api.get('/srsx/flights');
    }

    getFlight(id) {
        return api.get('/srsx/flight/' + id);
    }
}

export default new FlightService();