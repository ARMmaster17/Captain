import React from 'react';
import 'regenerator-runtime/runtime';
import axios from 'axios';
import {withRouter} from "react-router-dom";

class FlightForm extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            Name: ''
        };
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleSubmit(event) {
        event.preventDefault();
        let newFlightPayload = {
            Name: this.state.Name,
            AirspaceID: this.props.airspaceID
        };
        let requestOptions = {
            method: 'POST',
            headers: { 'ContentType': 'application/json' },
            body: JSON.stringify(newFlightPayload),
        };
        fetch('http://172.27.67.219:5000/flight', requestOptions)
            .then(response => response.json())
            .then(data => console.log(data))
            .catch(console.log);
    }

    render() {
        return (
            <form onSubmit={this.handleSubmit}>
                <label>
                    Name:
                    <input name={"inputName"} type={"text"} onChange={event => this.setState({Name: event.target.value})} />
                </label>
                <input type={"submit"} value={"Submit"} onSubmit={this.handleSubmit} />
            </form>
        )
    }
}

export default withRouter(FlightForm);