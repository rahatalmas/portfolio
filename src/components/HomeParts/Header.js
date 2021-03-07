import React from 'react';
import {BrowserRouter,Switch,Route,Link} from 'react-router-dom';
import './header.css';
const Header = () =>{
 
     return(

       <div className="headrow">
           <h3>  	&lt; &nbsp; LIFE OF ALMAS &nbsp; &gt;</h3>   
              <ul className="nav">
                  <Link to="/"> <li className="navli">HOME </li> </Link>
                  <Link to="/about"> <li className="navli"> ABOUT </li> </Link>
                  <Link to="/contact"> <li className="navli"> CONTACT </li> </Link>
              </ul>    
       </div>

     );

}

export default Header;