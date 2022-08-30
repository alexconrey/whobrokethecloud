// import { useEffect } from "react/cjs/react.production.min";
import { useState, useEffect } from 'react';
import Blame from './Blame';

function Home() {
    return (
        <Blame />
    )
    // const [error, setError] = useState(null);
    // const [isLoaded, setIsLoaded] = useState(false);
    // const [items, setItems] = useState([]);

    // var outageDict = new Object();
    // var currentOutageCount = new Object();

    // useEffect(() => {
    //     fetch("http://localhost/api/outages")
    //     .then(res => res.json())
    //     .then(
    //         (result) => {
    //             setIsLoaded(true);
    //             setItems(result);
    //         },
    //         (error) => {
    //             setIsLoaded(true);
    //             setError(error);
    //         }
    //     )
    // }, [])

    // if (error) {
    //     return <div>Error: {error.message}</div>
    // } else if (!isLoaded) {
    //     return <div>Loading...</div>;
    // } else {
    //     for (const [key, value] of Object.entries(items)) {
    //         outageDict[key] = value
    //     }

    //     for (const [key, value] of Object.entries(outageDict)) {
    //         currentOutageCount[key] = value.length
    //     }

    //     return (
    //         Object.entries(currentOutageCount)
    //             .map(([key, val]) => {
    //                 return (
    //                     <>
    //                         <li>{key}</li>
    //                         <li>{val}</li>
    //                     </>
    //                 )
    //                 console.log(val)
    //         })
    //     )
    // }
}

export default Home;