import React from 'react';
import IpLocationFormPage from "./features/ipsLookup/IpLocationFormPage";

function App() {
  return (
    <div data-testid="my-app" className="App">
        <IpLocationFormPage initIpLocations={[]}/>
    </div>
  );
}

export default App;
