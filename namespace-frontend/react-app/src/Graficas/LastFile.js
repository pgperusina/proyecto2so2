import React, { useState, useEffect } from "react";
import { Card } from "react-bootstrap";

export function LastFile(props) {

    const [data, setData] = useState(null)

    useEffect(() => {
        fetch(`${props.url}/last`)
            .then(res => res.json())
            .then(caso => {
                setData(JSON.parse(caso.value))
            })
        setInterval(() => {
            fetch(`${props.url}/last`)
                .then(res => res.json())
                .then(caso => {
                    setData(JSON.parse(caso.value))
                })
        }, 5000);
    }, []);

    return (
        <Card>
            <Card.Header>
                <Card.Title>Ultimo caso agregado</Card.Title>
            </Card.Header>
            <Card.Body>
                {
                    data !== null ?
                        <div>
                            <p><b>Nombre</b> <br /> {data.name}</p>
                            <p><b>Ubicación</b><br /> {data.location}</p>
                            <p><b>Edad</b><br /> {data.age}</p>
                            <p><b>Tipo de infección</b><br /> {data.infectedtype}</p>
                            <p><b>Estado</b><br /> {data.state}</p>
                        </div>
                        : null
                }
            </Card.Body>
        </Card>
    )
}