import { useState, useEffect } from 'react';
import Moment from 'moment';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Typography from '@mui/material/Typography';


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
            {/* <td>{props.starttime}</td> */}
            <td>{props.modifiedtime}</td>
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
                    // starttime={outage.StartTime}
                    modifiedtime={outage.ModifiedTime}
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
                    {/* <th>Start Time</th> */}
                    <th>Last Updated</th>
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

    Moment.locale('en')

    var outageDict = {};
    var currentOutageCount = {};
    var outageArr = [];

    useEffect(() => {
        const host = window.location.protocol + "//" + window.location.host;
        fetch(host+"/api/outages")
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
            return (
                <Container>
                    <Typography variant="h2" align="center">
                        No outages, surprisingly!
                    </Typography>
                    <Typography align="center">
                        Think we're wrong? <a href="mailto:dev@null.com">Let us know!</a>
                    </Typography>
                </Container>
            )
        } else if (affectedProviders === 1) {
            const outages = Object.entries(outageArr)[0][1]
            const outage = outages[0]
            // Tuesday, August 30 @ 11:32:35 AM
            const startTimeFormatted = Moment(outage.StartTime).format('dddd, MMMM DD YYYY @ hh:mm:ss A')
            return (
                <Box sx={{ my: 4 }}>
                    <h1>We're blaming {Object.keys(currentOutageCount)[0].capitalize()}</h1>
                    <div>
                        <p>On {startTimeFormatted}, {outage.Provider.capitalize()} reported an issue.</p>
                        <RenderOutageTable outages={outageArr} />
                    </div>
                </Box>
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