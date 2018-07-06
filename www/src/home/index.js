import React from 'react'

import { Link } from 'react-router-dom'

// import 'bootstrap';
// import 'bootstrap/dist/css/bootstrap.min.css';
import '../bootstrap-overrides.css';
import './index.css'

const Home = () => {
    return <div className="">
        <div className="bg-dark collapse" id="navbarHeader">
        <div className="container">
          <div className="row">
            <div className="col-sm-8 col-md-7 py-4">
              <h4 className="text-white">About</h4>
              <p className="text-muted">My name is Liam, I've been hacking since I was 11, and here's the first release of a product I've been working on.</p>
            </div>
            <div className="col-sm-4 offset-md-1 py-4">
              <h4 className="text-white">Contact</h4>
              <ul className="list-unstyled">
                <li><a href="https://twitter.com/liamzebedee" className="text-white">Follow on Twitter</a></li>
                <li><a href="mailto:basemap@liamz.co" className="text-white">Email</a></li>
              </ul>
            </div>
          </div>
        </div>
      </div>

      <div className="navbar navbar-dark bg-dark box-shadow">
        <div className="container d-flex justify-content-between">
          <a href="#" className="navbar-brand d-flex align-items-center">
            <strong>Basemap</strong>
          </a>
          <button className="navbar-toggler collapsed" type="button" data-toggle="collapse" data-target="#navbarHeader" aria-controls="navbarHeader" aria-expanded="false" aria-label="Toggle navigation">
            <span className="navbar-toggler-icon"></span>
          </button>
        </div>
      </div>

      <main role="main">

      <section className="jumbotron text-center">
        <div className="container">
          <h1 className="jumbotron-heading">Understand Go codebases at scale</h1>
          <div className="lead text-muted">
          <ul>
                <li>Quickly scan from a birds-eye view where functionality is located</li>
                <li>Understand relationships between types (usages, call graphs) without searching through files and directories</li>
                <li>Generate architectural diagrams, from the code itself</li>
            </ul>
          </div>
          <p>
            <Link to="/repo/?name=github.com/twitchyliquid64/subnet/subnet" className="btn btn-primary my-2">Show example</Link>
            &nbsp;
            <a href="#" className="btn btn-secondary my-2">Generate from GitHub repo</a>
          </p>
        </div>
      </section>

    </main>
        

        {/* <div className="input-group">
            <input type="text" className="form-control" placeholder="Enter a Go git repository e.g. https://github.com/liamzebedee/gitmonitor"/>
            <div className="input-group-append">
                <button className="btn btn-outline-secondary" type="button">Submit</button>
            </div>
        </div> */}

    </div>
}
export default Home;