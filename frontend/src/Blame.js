import { useState, useEffect } from 'react';

// eslint-disable-next-line
Object.defineProperty(String.prototype, 'capitalize', {
    value: function() {
      return this.charAt(0).toUpperCase() + this.slice(1);
    },
    enumerable: false
  });

function RenderOutageRow(props) {
    return (
        <tr key="{props.provider}-{props.service}">
            <td>{props.provider.capitalize()}</td>
            <td>{props.service}</td>
            <td>{props.starttime}</td>
            <td>{props.description}</td>
        </tr>
    )
}
function RenderOutageTable(props) {
    var outages = [];
    props.outages.map((outageArray) => {
        outageArray.map((outage) => {
            outages.push(
                <RenderOutageRow 
                    provider={outage.Provider} 
                    service={outage.Service} 
                    starttime={outage.StartTime} 
                    description={outage.Description} 
                />
            )
        })
    })
    return (
        <table>
            <thead>
                <tr>
                    <th>Provider</th>
                    <th>Service</th>
                    <th>Start Time</th>
                    <th>Description</th>
                </tr>
            </thead>
            <tbody>
                {outages}
            </tbody>
        </table>

    )
}

export default function Blame() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [items, setItems] = useState([]);

    var outageDict = {};
    var currentOutageCount = {};
    var outageArr = [];

    useEffect(() => {
        fetch("http://localhost/api/outages")
        .then(res => res.json())
        .then(
            (result) => {
                setIsLoaded(true);
                setItems(result);
            },
            (error) => {
                setIsLoaded(true);
                setError(error);
            }
        )
    }, [])

    if (error) {
        return <div>Error: {error.message}</div>
    } else if (!isLoaded) {
        return <div>Loading...</div>;
    } else {
        for (const [key, value] of Object.entries(items)) {
            if (! Array.isArray(outageDict[key])) {
                outageDict[key] = [];
            }
            outageDict[key].push(value)
            outageArr.push(value)
        }

        for (const [key, value] of Object.entries(outageDict)) {
            currentOutageCount[key] = value.length
        }

        const affectedProviders = Object.keys(currentOutageCount).length;

        if (affectedProviders === 0) {
            return <div>No outages, surprisingly!</div>
        } else if (affectedProviders == 1) {
            return (
                <div>
                    <h1>We're blaming {Object.keys(currentOutageCount)[0].capitalize()}</h1>
                    <div>
                        <RenderOutageTable outages={outageArr} />
                    </div>
                </div>
            )
        } else {
            return (
                <div>
                    <h1>Well, the internet is on fire.</h1>
                    <div>
                        <div>Current Outages</div>
                        <RenderOutageTable outages={outageArr} />
                    </div> 
                </div>
            )
        }
    }
}