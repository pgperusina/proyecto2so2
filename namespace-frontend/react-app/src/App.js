import logo from './logo.svg';
import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import { Container, Row, Col } from "react-bootstrap";
import { Top } from './Graficas/Top';
import { LastFile } from './Graficas/LastFile';
import { Pie } from './Graficas/Pie';
import { Bar } from './Graficas/Bar';

function App() {
  const url = "http://34.121.169.231:50501/api"
  return (
    <div className="App">
      <Container className="mb-3 mt-3">
        <Row>
          <Col md={12} className="text-center"><h1><b>Reporte de Casos Coronavirus</b></h1></Col>
          <Col md={3} >
            <Top url={url} />
          </Col>
          <Col md={9}>
            <Pie type="pie" name={"barra1"} url={url} />
          </Col>
          <Col md={3}>
            <LastFile url={url} />
          </Col>
          <Col md={9}>
            <Bar name={"barra2"} url={url} />
          </Col>
        </Row>
      </Container>
    </div>
  );
}

export default App;
