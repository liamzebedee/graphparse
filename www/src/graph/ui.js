import React from 'react';

import {
    searchNodes,
    selectNodeFromSearch,
    selectNodeByLabel
} from './actions'

import { connect } from 'react-redux'

import './ui.css';
import nodeColor from './colours';

class GraphControls extends React.Component {
    constructor() {
        super()
    }

    componentDidMount() {
        this.props.firstLoad()
    }

    render() {
        let { nodeTypes, q, matches, searchNodes, selectNode, clickedNode } = this.props;
        
        return <div className="info-view">
            {/* <div className='node-types'>
                {nodeTypes.map((typ, i) => {
                    return <span style={{ backgroundColor: nodeColor(i), color: 'white' }}>{typ}</span> 
                })}
            </div> */}

            <div>
                <input type='text' onChange={(ev) => searchNodes(ev.target.value)} value={q}/>

                <div className='results'>
                    { matches.length > 0 ? matches.map((node, i) => {
                        return <NodeMatch key={i} onClick={() => selectNode(node.id)} {...node}/>
                    }) : 'none' }
                </div>
            </div>

            <div className='debug'>
                <pre>{ clickedNode ? clickedNode.debugInfo : null }</pre>
            </div>
        </div>
    }
}


const NodeMatch = ({ onClick, label }) => {
    return <div onClick={onClick}>{label}</div>
}


const mapStateToProps = state => {
    return {
        ...state.graph.search,
        nodeTypes: state.graph.nodeTypes,
        clickedNode: state.graph.clickedNode ? state.graph.nodeLookup[state.graph.clickedNode] : null
    }
}

const mapDispatchToProps = dispatch => {
    return {
        firstLoad:   () => {
            dispatch(selectNodeByLabel("parse.go"))
            dispatch(searchNodes(".go"))
        },
        searchNodes: (q) => dispatch(searchNodes(q)),
        selectNode:  (id) => dispatch(selectNodeFromSearch(id))
    }
}

const GraphControlsView = connect(
  mapStateToProps,
  mapDispatchToProps
)(GraphControls)
 
export default GraphControlsView;