import {Route, Switch, withRouter} from "react-router-dom";
import React from "react";
import Flight from './Flight'
import FlightForm from './FlightForm'

class Airspace extends React.Component {
    state = {
        flights: []
    }

    componentDidMount() {
        // TODO: Change the URL.
        fetch(`http://172.27.67.219:5000/airspace/${this.props.match.params.airspaceId}/flights`)
            .then(res => res.json())
            .then((data) => {
                this.setState({
                    flights: data
                })
            })
            .catch(console.log)
    }

    render() {
        return (
            <div>
                <h3>Airspace: {this.props.match.params.airspaceId}</h3>
                <ul>
                    {this.state.flights?.map((flight) => (
                        <li key={flight.Name}>
                            <p>{flight.Name}</p>
                        </li>
                    ))}
                </ul>
                <Switch>
                    <Route path={`${this.props.match.path}/:airspaceId/:flightId`}>
                        <Flight />
                    </Route>
                    <Route path={this.props.match.path}>
                        <FlightForm airspaceID={this.props.match.params.airspaceId} />
                    </Route>
                </Switch>
            </div>
        )
    }
}

export default withRouter(Airspace);