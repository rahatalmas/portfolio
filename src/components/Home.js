import React from 'react';
import {useState,useEffect} from 'react';
import Body from './HomeParts/body';
import Header from './HomeParts/Header';
import {BrowserRouter,Switch,Route,Link} from 'react-router-dom';
import About from './HomeParts/About';
import Contact from './HomeParts/Contact';

const Home = () => {
    return (
     <div>
         <BrowserRouter>
         <Header/>
          <Switch>
          <Route exact path="/" component={Body}/>
          <Route path="/about" component={About}/>
          <Route path="/contact" component={Contact}/>
          </Switch> 
       </BrowserRouter>

     </div>
    );
};

export default Home;