import React from 'react';
import './aboutonecs.css';
const Aboutwo =(props)=>{
     return(
        <div>
            <div className="rowabout">
            <div className="cabout">
              <div>
              <h3>{props.heading}</h3>
               <p>{props.hello}</p>
               <p>{props.text}</p>
              </div> 
           </div>    
           <div className="cabout" className="c1about">
               <img src={props.imgsrc} className="imageabout"/>
           </div> 
       </div>
        </div>
       
     );
}

export default Aboutwo;