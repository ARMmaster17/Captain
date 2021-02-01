import React from 'react'
import ReactDOM from 'react-dom'
import PropTypes from 'prop-types'

class MachineList extends React.Component {
    state = {
        machines: []
    }

    componentDidMount() {
        this.fetchMachineData()
        setInterval(this.fetchMachineData, 5000)
    }
    fetchMachineData = () => {
        fetch('/machines.json')
            .then(res => res.json())
            .then((data) => {
                this.setState({ machines: data })
                console.log(this.state.machines)
            })
            .catch(console.log)
    }
    render() {
        return (
            <table className="table table-bordered">
                <thead className="thead-dark">
                    <tr>
                        <th>Hostname</th>
                        <th>Ip address</th>
                        <th>Cpu</th>
                        <th>Ram</th>
                        <th>Disk</th>
                        <th colSpan="3"></th>
                    </tr>
                </thead>
                <tbody>
                    {this.state.machines.map((machine) => (
                        <tr key={machine.id}>
                            <td>{machine.hostname}</td>
                            <td>{machine.ip_address}</td>
                            <td>{machine.cpu}</td>
                            <td>{machine.ram}</td>
                            <td>{machine.disk}</td>
                            <td><a rel="nofollow" data-method="delete" href={"/machines/" + machine.id} onClick="return false">Destroy</a></td>
                        </tr>
                    ))}
                </tbody>
            </table>
        )
    }
}

export default MachineList