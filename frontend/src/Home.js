// import { useEffect } from "react/cjs/react.production.min";
import { useState, useEffect } from 'react';

function Home() {
    const [error, setError] = useState(null);
    const [isLoaded, setIsLoaded] = useState(false);
    const [items, setItems] = useState([]);

    var outageDict = new Object();
    var currentOutages = new Object();

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
            outageDict[key] = value
        }

        for (const [key, value] of Object.entries(outageDict)) {
            currentOutages[key] = value.length
        }

        return (
            Object.entries(currentOutages)
                .map(([key, val]) => {
                    return (
                        <>
                            <li>{key}</li>
                            <li>{val}</li>
                        </>
                    )
                    console.log(val)
            })
        )

        // return (
        //     // {Object.entries(outageDict)
        //     //     .map(([key, val]) => {
        //     //         console.log(val)
        //     //         return (
        //     //             <ul id={key}>
        //     //                 {val.map(element => {
        //     //                     return <li>{element.Service}</li> 
        //     //                 })}
        //     //             </ul>
        //     //         )              
        //     //     })}
        //     Object.entries(outageDict)
        //     .map(([key, val]) => {
        //         console.log(val)
        //         return (
        //             <ul id={key}>
        //                 {Object.entries(val)
        //                 .map(([k, v])) => {
        //                     return ( <li>{v}</li> )
        //                 }}
        //                 <li id={key}>{val.Service}</li>
        //             </ul>
        //         )              
        //     })
        // )
        // console.log(outageDict)
        // console.log(items)
        // return (
        //     <ul>
        //         {items.map((element, idx) => (
        //             {element}
        //         ))}
        //     </ul>
        // );
    }
}

export default Home;