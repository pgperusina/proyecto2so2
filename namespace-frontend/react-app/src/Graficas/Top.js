import React, { useState, useEffect } from "react";
import { Card } from "react-bootstrap";

export function Top(props) {

    const [data, setData] = useState([])

    useEffect(() => {
        console.log(`getting info from : ${props.top}/top`);
        console.log(props)
        fetch(`${props.url}/top`)
            .then(res => res.json())
            .then(casos => {
                setData(casos.casos)
            })
        setInterval(() => {
            fetch(`${props.url}/top`)
                .then(res => res.json())
                .then(casos => {
                    setData(casos.casos)
                })
        }, 5000);
    }, []);

    return (
        <Card className="text-white bg-info mb-2">
            <Card.Header>
                <Card.Title>Top 3 de Departamentos con Coronavirus</Card.Title>
            </Card.Header>
            <Card.Body>
                {
                    data.map((caso, index) => {
                        return (
                            <Card className="text-white bg-success mb-2" key={index}>
                                <Card.Header>{caso._id}</Card.Header>
                                <Card.Body>
                                    {caso.count} Casos
                                </Card.Body>
                            </Card>
                        )
                    })
                }
                {/* <Card className="text-white bg-success mb-2">
                    <Card.Header>Peten</Card.Header>
                    <Card.Body>
                        30 Casos
                    </Card.Body>
                </Card>
                <Card className="text-white bg-danger mb-2">
                    <Card.Header>El Progreso</Card.Header>
                    <Card.Body>
                        15 Casos
                    </Card.Body>
                </Card> */}
            </Card.Body>
        </Card>
    )
}