import {fireEvent, render, screen, waitFor} from '@testing-library/react';
import App from './App';
import {Provider} from 'react-redux';
import {store} from './app/store';
import React, {useState} from "react";
import {IpLocation} from "./features/ipsLookup/type";
import IpLocationFormPage from "./features/ipsLookup/IpLocationFormPage";

describe('renders App component', () => {
    it('should display normally', () => {
        render(
            <Provider store={store}>
                <App/>
            </Provider>
        );

        expect(screen.getByTestId('my-app')).toBeInTheDocument();
    });

    it('should display error when input empty IPs', async () => {
        render(
            <Provider store={store}>
                <App/>
            </Provider>
        );

        const submitButton = screen.getByRole('button', {name: 'Submit'});
        const fetchSpy = jest.spyOn(window, 'fetch');

        fireEvent.click(submitButton);
        expect(fetchSpy).toHaveBeenCalled();

        await waitFor(() => expect(screen.getByText('Please input correct IPs')).toBeInTheDocument());
    });

    it('should display error when input wrong IPs format', async () => {
        render(
            <Provider store={store}>
                <App/>
            </Provider>
        );

        const submitButton = screen.getByRole('button', {name: 'Submit'});
        const fetchSpy = jest.spyOn(window, 'fetch');

        fireEvent.change(screen.getByTestId('ips-input-id'), {target: {value: '1.1.1.1; 2.2.2.2;'}})
        fireEvent.click(submitButton);
        expect(fetchSpy).toHaveBeenCalled();

        await waitFor(() => expect(screen.getByText('Please input correct IPs')).toBeInTheDocument());
    });

    it('should display error when input include any IPs format', async () => {
        render(
            <Provider store={store}>
                <App/>
            </Provider>
        );

        const submitButton = screen.getByRole('button', {name: 'Submit'});
        const fetchSpy = jest.spyOn(window, 'fetch');

        fireEvent.change(screen.getByTestId('ips-input-id'), {target: {value: '1.1.1.1, 2.2.2.256'}})
        fireEvent.click(submitButton);
        expect(fetchSpy).toHaveBeenCalled();

        await waitFor(() => expect(screen.getByText('Please input correct IPs')).toBeInTheDocument());
    });

    it('should display result normally', async () => {
        const testIpLocations: IpLocation[] = [
            { ip: '1.1.1.1', city: 'Toronto', timezone: 'Canada', accuracy_radius: 0, postal_code: 'K1K' },
        ];

        render(
                <IpLocationFormPage initIpLocations={testIpLocations}/>
        );
        expect(screen.getByText('Canada')).toBeInTheDocument();

    });
});
