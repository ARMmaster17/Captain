import './App.css';

import React from 'react';
import {
    BrowserRouter as Router,
    Switch,
    Route,
    Link,
    useParams,
    useRouteMatch
} from "react-router-dom";
import Airspaces from './Airspaces'

class App extends React.Component {
    render() {
        return (
            <Router>
                <div>
                    <nav>
                        <ul>
                            <li>
                                <Link to={"/"}>Home</Link>
                            </li>
                            <li>
                                <Link to={"/airspaces"}>Airspaces</Link>
                            </li>
                        </ul>
                    </nav>
                    <Switch>
                        <Route path={"/airspaces"}>
                            <Airspaces />
                        </Route>
                        <Route path={"/"}>
                            <Home />
                        </Route>
                    </Switch>
                </div>
            </Router>
        );
    }
}

function Home() {
    return <h2>Home</h2>;
}

export default App;

