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
    clickNode
} from './actions'
import {
    hexToRgb
} from '../util'
import nodeColor from './colours';
import styles from './graph.css';

const nodeType = (str) => graphJSON.nodeTypes.indexOf(str);

class D3Graph extends React.Component {
    constructor() {
        super()
        this.graphDOT = null;
    }

    state = {
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
    }

    addZoom = () => {
        var zoomHandler = d3.zoom()
        .on("zoom", () => {
            this.setState({ zoom: d3.event.transform })
        });
  
        zoomHandler(d3.select(this.svg));
    }

    static getDerivedStateFromProps(nextProps, prevState) {
        let { nodeLookup, edges } = nextProps;

        const showOnlyInterestedNodes = (edge) => {
            return _.contains(nextProps.interested, edge.source)
        }
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
            if(_.contains(seenEdges, id)) return null;
            seenEdges.push(id)
            return true
        }
        edges = edges
        .filter(showOnlyInterestedNodes)
        .filter(collapseRedundantEdges)
        .filter(removeDuplicates)



        let nodesToInclude = edges.map(e => [e.source, e.target]).reduce((a, b) => a.concat(b), []);
        let nodes = nextProps.nodes.filter(node => {
            return _.contains(nodesToInclude, node.id)
        })

        if(nodes.length < 1 || edges.length < 1) return {}
        

        let graphDOT = generateGraphDOT(nodes, edges)

        let layout = generateLayout(graphDOT, nextProps.nodeLookup)
        
        return {
            graphDOT,
            ...layout,
        }
    }

    render() {
        let zoom = this.state.zoom;
        let seenEdges = [];
        

        return <div>
            <svg width='100%' height='100%' ref={(svg) => this.svg = svg}>
                <g class='everything' style={{
                    transform: `translate3d(${zoom.x}px, ${zoom.y}px, 0px) scale(${zoom.k})`
                }}>
                    <g class='nodes'>
                        {this.state.nodes.map(node => {
                            return <Node 
                                key={node.id} clickNode={this.props.clickNode} 
                                interesting={_.contains(this.props.interested, node.id)}
                                {...node}/>
                        })}
                    </g>

                    <g class='edges'>
                        {this.state.edges.map((edge, i) => {
                            if(_.contains(seenEdges, edge.id)) return false;
                            seenEdges.push(edge.id)
                            return <Edge key={edge.id} {...edge}/>
                        })}
                    </g>
                </g>
            </svg>
        </div>
    }
}

const Node = ({ id, interesting, cx, cy, _draw_, _ldraw_, variant, label, clickNode }) => {
    return <g 
        class='node'
        onMouseOver={() => hoverNode(id)}
        onClick={() => clickNode(id)}>
        <ellipse 
            stroke='#000000'
            cx={cx}
            cy={cy}
            rx={_draw_[1].rect[2]}
            ry={_draw_[1].rect[3]}
            fill={nodeColor(variant)}
            class={classNames({
                'interesting': interesting
            })}
        />
        <text 
            textAnchor='middle'
            x={_ldraw_[2].pt[0]}
            y={-_ldraw_[2].pt[1]}>
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


const toSvgPointSpace = point => [ point[0], -point[1] ];


const generateGraphDOT = (nodes, edges) => `
    digraph graphname {
        ${nodes.map(({ id, rank, label }) => {
            rank = 1;
            return `"${id}" [width=${rank}] [height=${rank}] [label="${label}"];`
        }).join('\n')}
        ${edges.map(({ source, target }) => `"${source}" -> "${target}";`).join('\n')}
    }
`

// // Passes DOT to Graphviz, generates layout of nodes and edges in JSON, merges with node data to be bound to D3
function generateLayout(graphDOT, nodeLookup) {
    let graphvizData = JSON.parse(Viz(graphDOT, { format: 'json' }));
    
    let nodes = graphvizData.objects.map(obj => {
        let pos = obj.pos.split(',').map(Number);
        let id = new Number(obj.name)

        return {
            cx: pos[0],
            cy: -pos[1],
            rx: obj._draw_[1].rect[2],
            ry: obj._draw_[1].rect[3],
            ...obj,
            id,
            ...nodeLookup[id]
        }
    })

    let edges = graphvizData.edges.map((edge, i) => {
        let points = edge._draw_[1].points.map(toSvgPointSpace);
        
        let { head, tail } = edge;
        function findNodeForGvid(id) {
            let obj = _.find(graphvizData.objects, obj => obj._gvid == id)
            if(!obj) throw new Error()
            return Number(obj.name)
        }

        let source = findNodeForGvid(head)
        let target = findNodeForGvid(tail)

        let arrowPts = edge._hdraw_[3].points.map(toSvgPointSpace)

        return {
            points,
            arrowPts,
            source,
            target,
            id: `${i}${source}${target}`
        }
    })

    return {
        nodes,
        edges,
    }
}


// function renderGraphviz(g) {
//     var graphvizSvg = Viz(graphDOT, { format: 'svg' });
//     g.html(graphvizSvg)
// }

// function focusOnPackageNode(svg, zoom, layout) {
//     let node = _.find(layout.nodes, node => {
//         return node.variant == graphJSON.nodeTypes.indexOf('RootPackage')
//     })
//     // let transform = to_bounding_box(getCenter(node.cx, node.cy), node.rx, node.ry, 0)
//     // g.transition().duration(200).call(zoom.transform, transform);
//     // zoom.translateTo(svg, node.cx / 2, node.cy / 2)
// }

// function contrastColour(r,g,b) {
//     let d = 0;
//     // Counting the perceptive luminance - human eye favors green color... 
//     let a = 1 - ( 0.299 * r + 0.587 * g + 0.114 * b)/255;
//     if (a < 0.5)
//         return "black";
//     else
//         return "white";
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
        }
    }
}

const D3GraphCtn = connect(
  mapStateToProps,
  mapDispatchToProps
)(D3Graph)
 
export default D3GraphCtn;