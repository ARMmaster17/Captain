import {Link, Route, Switch, withRouter} from "react-router-dom";
import React from "react";
import Airspace from './Airspace'

class Airspaces extends React.Component {
    state = {
        airspaces: []
    }

    componentDidMount() {
        // TODO: Change to... something else.
        fetch('http://172.27.67.219:5000/airspaces')
            .then(res => res.json())
            .then((data) => {
                this.setState({
                    airspaces: data
                })
            })
            .catch(console.log)
    }

    render() {
        return (
            <div>
                <h2>Airspaces</h2>
                <ul>
                    {this.state.airspaces.map((airspace) => (
                        <li key={airspace.NetName}>
                            <Link to={`${this.props.match.url}/${airspace.ID}`}>{airspace.HumanName}</Link>
                        </li>
                        ))}
                </ul>
                <Switch>
                    <Route path={`${this.props.match.path}/:airspaceId`}>
                        <Airspace />
                    </Route>
                    <Route path={this.props.match.path}>
                        <h3>Please select an airspace.</h3>
                    </Route>
                </Switch>
            </div>
        )
    }
}

export default withRouter(Airspaces);