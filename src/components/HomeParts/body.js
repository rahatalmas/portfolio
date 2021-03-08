import React from 'react';
import {useState,useEffect} from 'react';
import Gallary from './Gallary';
import "./body.css"
const Body = () => {
    return (
        <div>
     < div className="row">
          <div className="cl" class="t1">
                 <div className="t1text">
                     <h2>Hello World</h2>
                     <h3>welcome to life of almas</h3>
                     <br/>
                     <code>
                         int main()
                         <br/>
                         eat ; sleep ; code ; repeat ;
                     </code>
                 </div>
          </div>

          <div className="cl">
            <div className="card">
                <img className="cardimg" src="141136622_2738010949791311_8621849036606549127_o.jpg"/>
                 <div className="cardbio">
                     <h3>Rahat Almas</h3>
                     <br/>
                     <p>Student of Computer Science and Engineering</p>
                     <br/>
                 </div>
            </div>
          </div>

          
     </div>
     <div className="editor">
       <p>Favourite Captures</p>
       </div>
       <Gallary/>
     </div>
    );
};

export default Body;