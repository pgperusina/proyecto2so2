import React, { useEffect, useState } from 'react'
import Chart from "chart.js";
import { Card } from 'react-bootstrap';

export function Pie(props) {

    const [titles, setTitles] = useState([])
    const [data, setData] = useState([])
    const [recargar, setRecargar] = useState(true)
    const ctx = document.getElementById(props.name);

    useEffect(() => {
        fetch(`${props.url}/pipeline`)
            .then(res => res.json())
            .then(casos => {
                let d = []
                let t = []
                casos.casos.forEach(caso => {
                    t.push(caso._id)
                    d.push(caso.count)
                });
                setData(d)
                setTitles(t)
            })
        setInterval(() => {
            fetch(`${props.url}/pipeline`)
                .then(res => res.json())
                .then(casos => {
                    let d = []
                    let t = []
                    casos.casos.forEach(caso => {
                        t.push(caso._id)
                        d.push(caso.count)
                    });
                    setData(d)
                    setTitles(t)
                })
        }, 5000);
    }, []);

    if (data.length > 0 && titles.length > 0) {
        new Chart(ctx, {
            type: props.type,
            data: {
                labels: titles,
                datasets: [
                    {
                        label: "# of Votes",
                        data: data,
                        backgroundColor: [
                            'rgba(255, 99, 132, 0.4)',
                            'rgba(54, 162, 235, 0.4)',
                            'rgba(255, 206, 86, 0.4)',
                            'rgba(75, 192, 192, 0.4)',
                            'rgba(153, 102, 255, 0.4)',
                            'rgba(255, 159, 64, 0.4)'
                        ],
                        borderColor: [
                            'rgba(255, 99, 132, 1)',
                            'rgba(54, 162, 235, 1)',
                            'rgba(255, 206, 86, 1)',
                            'rgba(75, 192, 192, 1)',
                            'rgba(153, 102, 255, 1)',
                            'rgba(255, 159, 64, 1)'
                        ],
                        hoverBackgroundColor: 'rgba(0, 0, 0, 0.2)',
                        borderWidth: 1
                    }
                ]
            }
        });
    }

    return (
        <Card>
            <Card.Header>
                <Card.Title>Departamentos afectados</Card.Title>
            </Card.Header>
            <Card.Body>
                <canvas id={props.name} width="100%"></canvas>
            </Card.Body>
        </Card>
    )
}