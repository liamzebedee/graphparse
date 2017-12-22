import 'script-loader!./vendor/d3.v4.min.js';

import './ast.css'

import React from 'react';
import ReactDOM from 'react-dom';
import copy from 'copy-to-clipboard';


document.addEventListener('DOMContentLoaded', function() {
    ReactDOM.render(
        <ASTView/>,
        document.getElementById('mount')
    );
});

class ASTView extends React.Component {
    constructor() {
        super()
        this.state = {
            src: {
                code:""
            },
            output: "",
        }
        this.copyOutput = this.copyOutput.bind(this);
    }

    componentWillMount() {
        // Fetch initial source file.
        fetch('http://localhost:8081/src', {})
        .then(response => response.json())
        .then(data => {
            this.setState({
                src: data
            })
        })

        document.addEventListener('keydown', (ev) => {
            if(ev.shiftKey) {
                let selection = window.getSelection();
                if(selection) {
                    this.onCodeSelection(selection.anchorOffset, selection.focusOffset)
                }
                return false;
            }
        })
    }

    onCodeSelection(start, end) {
        start = this.state.src.pos + start
        end = this.state.src.pos + end
        
        fetch(`http://localhost:8081/src/from/${start}/to/${end}`)
        .then(response => response.json())
        .then(data => {
            this.setState({
                output: data.output
            })
        })
    }

    getPosForLine(i) {
        let lines = this.state.src.code.split('\n')
        let posAcc = 0;
        lines.map((line, j) => {
            if(j == i) return;
            posAcc += line.length;
        })
        return this.state.src.pos + posAcc;
    }
    
    copyOutput() {
        copy(this.state.output, {
            debug: true,
            message: 'Press #{key} to copy',
        });
    }

    render() {
        let lines = this.state.src.code.split('\n')

        return <div className='container-fluid'>
            <div id='code'>
                <pre ref={(ref) => this.codePre = ref}>
                    {this.state.src.code || <div className="alert alert-warning" role="alert">
                    Can't load code. Is graphparse running with the API server flag? <pre>graphparse -api</pre>
                    </div>}
                </pre>
            </div>
            <div id='sidebar'>
                <div className='wrap'>
                    {this.state.output 
                    ? <a href='#' onClick={this.copyOutput} className="btn btn-light">Copy</a>
                    : null }
                    
                    <pre>
                        {this.state.output 
                        ? this.state.output
                        : <div className="alert alert-info" role="alert">
                        Select code and press [shift] to see its AST representation
                        </div>}
                    </pre>
                </div>
            </div>
        </div>;
    }
}