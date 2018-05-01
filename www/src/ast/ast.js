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

function parseASTOutput(output) {
    let root = [];
    const indent = '  ';

    // get head of first line
    let lines = output.split('\n')
    let lineNumIndent = lines[0].split('0')[0].length;
    lines = lines.map(line => line.slice(lineNumIndent + indent)).join('\n')

    function parseBlock(blockStr) {
        let currentIndent = 1;
        let children = [];

        blockStr.split('\n').map((line, i) => {
            let indents = line.split('.  ');
            if(indents.length === 1) {
                return
            }
            if(indents.length + 1 === currentIndent) {
                if(!indents[indents.length] == '}')
                    children.push(i)
            }
        })

        return children
    }

    console.log(parseBlock(lines))
}

// <TreeView children={}/>

// children.map(<Treeview children={}/>)

class ASTView extends React.Component {
    constructor() {
        super()
        this.state = {
            src: {
                code:"",
                pos: -1,
            },
            output: "",
        }
        this.copyOutput = this.copyOutput.bind(this);
    }

    componentWillMount() {
        // function getCorrectPosForFile(code, origPos) {
        //     
        //     return origPos - code.indexOf("package")
        // }

        function trimCodeToStartPos(code, origPos) {
            // ast.File.Pos() is the position of "package" keyword
            return code.split('').splice(code.indexOf("package")).join('')
        }

        // Fetch initial source file.
        fetch('http://localhost:8081/src/', {})
        .then(response => response.json())
        .then(data => {
            this.setState({
                src: {
                    code: trimCodeToStartPos(data.code, data.pos),
                    pos: data.pos,
                    // pos: getCorrectPosForFile(data.code, data.pos)
                },
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
        start = this.state.src.pos + start;
        end = this.state.src.pos + end;
        
        fetch(`http://localhost:8081/src/ast-range?start=${start}&end=${end}` )
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
                    {this.state.src.code === undefined ? <div className="alert alert-warning" role="alert">
                    Can't load code. Is graphparse running with the API server flag? <pre>graphparse -api</pre>
                    </div> : null }
                    {this.state.src.code == "" ? <div className="alert alert-warning" role="alert">
                    Server responded with no code. Check package / file import path.
                    </div> : null }
                    {this.state.src.code}
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

