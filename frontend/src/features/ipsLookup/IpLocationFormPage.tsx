import React, {useEffect, useState} from 'react';
import {v4 as uuidv4} from 'uuid'
import {Alert, Button, Container, Form, Row, Table} from "react-bootstrap";
import {IpLocation} from "./type";


const IpLocationFormPage: React.FC<{initIpLocations:IpLocation[]}> = ({initIpLocations}) => {

    const ENDPOINT = 'http://localhost:8080';

    const [ipsInput, setIpsInput] = useState<string>('');
    const [ipLocations, setIpLocations] = useState<IpLocation[]>(initIpLocations);
    const [inputError, setInputError] = useState<string>('')
    const [submitted, setSubmitted] = useState(false);
    const [clientId] = useState<string>(uuidv4())

    const validateIps = (value: string) => {
        try {
            setInputError('');
            const ipsArray = value.split(',').map(ip => ip.trim());
            for (let ip of ipsArray) {
                if (!/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(ip)) {
                    setInputError('Please input correct IPs');
                    return false;
                }
            }
            return true;
        } catch (e) {
            setInputError('Please input correct IPs in CSV format');
            return false;
        }
    };


    const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        validateIps(ipsInput);
        setSubmitted(true)
        const ipsArray = ipsInput.split(',').map(ip => ip.trim());
        const uniqueIpsSet = new Set(ipsArray);
        const uniqueIpsArray = Array.from(uniqueIpsSet);
        const ipsJson = JSON.stringify({"ip": uniqueIpsArray});

        const url = ENDPOINT + '/ips/' + clientId;

        await fetch(url, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: ipsJson
        });
    }

    useEffect(() => {
        if (process.env.NODE_ENV !== 'test') {
            const eventSourceUrl = ENDPOINT + '/stream/' + clientId;
            const eventSource = new EventSource(eventSourceUrl);

            eventSource.addEventListener('message', (event) => {
                const response = event.data;
                console.log(response);
                setIpLocations(JSON.parse(response));
                setSubmitted(false);
            });

            return () => {
                eventSource.close();
            };
        }
    }, [clientId]);

    return (
        <Container>
            <Alert variant={"info"}>
                Check your ip locations
            </Alert>
            <Form onSubmit={handleSubmit}>
                <Form.Group className={"mb-3"}>
                    <Form.Label>Enter Your IPs separated by comma (ex: 0.0.0.0, 1.1.1.1)</Form.Label>
                    <Form.Control
                        data-testid="ips-input-id"
                        as="textarea"
                        rows={3}
                        placeholder="ex. (0.0.0.0, 1.1.1.1)"
                        value={ipsInput}
                        onChange={(e) => setIpsInput(e.target.value)}
                    />
                </Form.Group>
                {
                    inputError &&
                    <Alert variant={"danger"}>
                        {inputError}
                    </Alert>
                }
                {!submitted && <Button variant={"primary"} type={"submit"}>Submit</Button>}
            </Form>
            <Row className={"my-3"}></Row>
            <Container>
                <Table striped bordered hover>
                    <thead>
                        <tr>
                            <th>IP</th>
                            <th>City</th>
                            <th>Timezone</th>
                            <th>AccuracyRadius</th>
                            <th>PostalCode</th>
                        </tr>
                    </thead>
                    <tbody>
                    {
                        ipLocations && ipLocations.map((ipLocation: IpLocation) => (
                            <tr key={ipLocation.ip}>
                                <td>{ipLocation.ip}</td>
                                <td>{ipLocation.city}</td>
                                <td>{ipLocation.timezone}</td>
                                <td>{ipLocation.accuracy_radius}</td>
                                <td>{ipLocation.postal_code}</td>
                            </tr>
                        ))
                    }
                    </tbody>
                </Table>
            </Container>
        </Container>
    );
};

export default IpLocationFormPage;
