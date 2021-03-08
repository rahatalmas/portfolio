import React from 'react';
import './aboutonecs.css';
const Aboutone =(props)=>{
     return(
        <div>
            <div className="rowabout">
           <div className="cabout" className="c1about">
           <img src={props.imgsrc} className="imageabout"/>
           <div>
              <h4>skills :</h4><br/>
              <ul>
                 <li>c++</li>
                 <li>java</li>
                 <li>javascript</li>
                 <li>react</li>
                 <li>python</li>

              </ul>
           </div>
           </div> 
           <div className="cabout">
              <div>
              <h3>{props.heading}</h3>
               <p>{props.hello}</p>
               <p>{props.text}</p>
              </div> 
           </div>
       </div>
        </div>
       
     );
}

export default Aboutone;