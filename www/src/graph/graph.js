import React from 'react';
import 'script-loader!../vendor/d3.v4.min.js';
import Viz from 'viz.js';
import _ from 'underscore';
import classNames from 'classnames';
import { connect } from 'react-redux'
import shortcut from 'keyboard-shortcut';
import copy from 'copy-to-clipboard';

import graphDOT from 'raw-loader!../../graph.dot';
import graphJSON from '../../graph.json';
import {
    hoverNode,
    clickNode,
    setGrabbing
} from './actions'
import {
    hexToRgb
} from '../util'
import nodeColor from './colours';

import './graph.css';

const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

class D3Graph extends React.Component {
    constructor() {
        super()
        this.graphDOT = null;
    }

    state = {
        grabbing: false,

        zoom: {
            transform: {
                x: 0,
                y: 0,
                z: 0
            }
        },
        nodes: [],
        edges: [],
        graphDOT: "",
    }

    componentDidMount() {
        this.addZoom()

        shortcut('ctrl c', {}, () => {
            copy(this.state.graphDOT);
        })

        // this.svg.addEventListener('mousedown', () => {          
        //     this.setState({ grabbing: true })
        // })
        // this.svg.addEventListener('mouseup', () => {          
        //     this.setState({ grabbing: false })
        // })
    }

    addZoom = () => {
        var zoomHandler = d3.zoom()
        .on("zoom", () => {
            this.setState({ zoom: d3.event.transform })
        });
  
        zoomHandler(d3.select(this.svg));
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        let { 
            interested, nodeLookup, 
            nodes, edges,
        } = nextProps;

        [nodes, edges] = filterNodesAndEdges({ nodes, edges, interested, nodeLookup });
        [nodes, edges] = generateLayout({ nodes, edges, nodeLookup });
        [nodes, edges] = postFilterNodesAndEdges({ nodes, edges, interested, nodeLookup });
        
        return {
            graphDOT,
            nodes,
            edges,
        }
    }

    render() {
        let zoom = this.state.zoom;

        return <svg
                className={classNames({
                    grabbing: this.state.grabbing
                })}
                // onTouchStart={window.alert}
                // onMouseMove={window.alert}
                // onTouchMove={window.alert}
                // onMouseUp={window.alert}
                // onTouchEnd={window.alert}
                ref={(ref) => this.svg = ref}
                >
            <defs>
                <filter id="shadow" x="0" y="0" width="200%" height="200%">
                <feOffset result="offOut" in="SourceAlpha" dx="10" dy="10" />
                <feGaussianBlur result="blurOut" in="offOut" stdDeviation="3" />
                <feBlend in="SourceGraphic" in2="blurOut" mode="normal" />
                </filter>
            </defs>

            <g 
                style={{
                    transform: `translate3d(${zoom.x}px, ${zoom.y}px, 0px) scale(${zoom.k})`
                }} 
                ref={ref => this.everything = ref}
                >

                <g>
                {this.state.nodes.map(node => {
                    return <Node 
                        key={node.id} clickNode={this.props.clickNode} 
                        interesting={_.contains(this.props.interested, node.id)}
                        {...node}/>
                })}
                </g>

                <g>
                {this.state.edges.map((edge, i) => {
                    return <Edge key={edge.id} {...edge}/>
                })}
                </g>
            </g>
        </svg>
    }
}

const Node = ({ id, interesting, cx, cy, _draw_, variant, label, clickNode }) => {
    return <g 
        class='node'
        onMouseOver={() => hoverNode(id)}
        onClick={() => clickNode(id)}
        transform={`translate(${cx}, ${cy})`}
        >
        <ellipse 
            stroke='#000000'
            rx={_draw_[1].rect[2]}
            ry={_draw_[1].rect[3]}
            fill={nodeColor(variant)}
            class={classNames({
                'interesting': interesting
            })}
        >
        </ellipse>
        <text 
            textAnchor='middle'
            x="0" y="0" alignment-baseline="middle" font-size="12" stroke-width="0" stroke="#000" text-anchor="middle"
            >
            {label}
        </text>
    </g>
}

const Edge = ({ points, arrowPts }) => {
    let computeD = () => {
        return points.map((point, i) => {
            if(i == 0) return `M${point.join(',')}C`;
            return `${point.join(',')} `;
        }).join('')
    }

    return <g>
        <path 
            fill='none'
            stroke='#000000'
            d={computeD()}/>
        <polygon
            fill="#000000"
            stroke="#000000"
            points={`${arrowPts.join(' ')} ${arrowPts[0]}`}/>
    </g>
}

const toSvgPointSpace = point => [ point[0], point[1] ];


const generateGraphDOT = (nodes, edges) => `
    digraph graphname {
        ${nodes.map(({ id, rank, label }) => {
            rank = 1;
            return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"];`
        }).join('\n')}
        ${edges.map(({ source, target, id }) => `"${target}" -> "${source}" [id=${id}];`).join('\n')}
    }
`


// Passes DOT to Graphviz, generates layout of nodes and edges in JSON
function generateLayout({ nodes, edges, nodeLookup }) {
    if(nodes.length < 1 || edges.length < 1) return [[], []];

    let graphDOT = generateGraphDOT(nodes, edges)

    let graphvizData = JSON.parse(Viz(graphDOT, { format: 'json' }));
    
    nodes = graphvizData.objects.map(obj => {
        let pos = obj.pos.split(',').map(Number);
        let id = new Number(obj.name)

        return {
            cx: pos[0],
            cy: pos[1],
            rx: obj._draw_[1].rect[2],
            ry: obj._draw_[1].rect[3],
            ...obj,

            id,
            ...nodeLookup[id]
        }
    })

    edges = graphvizData.edges.map((edge, i) => {
        let points = edge._draw_[1].points.map(toSvgPointSpace);
        
        // function findNodeForGvid(id) {
        //     let obj = _.find(graphvizData.objects, obj => obj._gvid == id)
        //     if(!obj) throw new Error()
        //     return Number(obj.name)
        // }

        // let { head, tail } = edge;
        // let source = findNodeForGvid(head)
        // let target = findNodeForGvid(tail)

        let edgeData = _.find(edges, edge_ => edge_.id == edge.id)

        let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace)

        return {
            points,
            arrowPts,
            ...edgeData
        }
    })

    return [ nodes, edges ];
}

function filterNodesAndEdges({ nodes, edges, interested, nodeLookup }) {
    const showOnlyInterestedNodes = (edge) => {
        return _.contains(interested, edge.source)
    }

    edges = edges
    .filter(showOnlyInterestedNodes)

    let nodesToInclude = edges.map(e => [e.source, e.target]).reduce((a, b) => a.concat(b), []);
    nodes = nodes.filter(node => {
        return _.contains(nodesToInclude, node.id)
    })

    return [ nodes, edges ]
}

function postFilterNodesAndEdges({ nodes, edges, interested, nodeLookup }) {
    const collapseRedundantEdges = (edge) => {
        let { source, target } = edge;
        let [ a, b ] = [ nodeLookup[source], nodeLookup[target] ]
        if(a.variant == nodeType('Struct')) {
           return !_.find(edges, (edge) => {
               return target === edge.target && source != edge.source
           })
        }
        return true;
    }

    let seenEdges = [];
    function removeDuplicates(edge) {
        let id = `${edge.source}${edge.target}`;
        if(_.contains(seenEdges, id)) return false;
        seenEdges.push(id)
        return true
    }

    edges = edges
    .filter(removeDuplicates)
    .filter(collapseRedundantEdges)

    return [ nodes, edges ]
}

// function zoomToBoundingBox(node) {
//     var bounds = path.bounds(d),
//       dx = bounds[1][0] - bounds[0][0],
//       dy = bounds[1][1] - bounds[0][1],
//       x = (bounds[0][0] + bounds[1][0]) / 2,
//       y = (bounds[0][1] + bounds[1][1]) / 2,
//       scale = Math.max(1, Math.min(8, 0.9 / Math.max(dx / width, dy / height))),
//       translate = [width / 2 - scale * x, height / 2 - scale * y];

//     let transform = d3.zoomIdentity.translate(translate[0],translate[1]).scale(scale);
//     return transform;
// }


const mapStateToProps = state => {
    return {
        ...state.graph,
    }
}

const mapDispatchToProps = dispatch => {
    return {
        clickNode: id => {
            dispatch(clickNode(id))
        },
        hoverNode: id => {
            dispatch(hoverNode(id))
        },
        setGrabbing: (grabbing) => dispatch(setGrabbing(grabbing)),
    }
}

const D3GraphCtn = connect(
  mapStateToProps,
  mapDispatchToProps
)(D3Graph)
 
export default D3GraphCtn;