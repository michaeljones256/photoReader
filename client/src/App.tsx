import React from 'react';
import logo from './logo.svg';
import MyCounter from './components/MyCounter'
import './App.css';

function App() {
  return (
    <div className="App">
      <MyCounter count={1}/>
    </div>
  );
}

export default App;
