import React from 'react';

import {
    searchNodes,
    loadInitialFileForTesting,
    selectNodeByLabel,
    selectNodeFromSearch,
} from './actions'

import { connect } from 'react-redux'

import './ui.css';
import nodeColor, { getVariantName } from './colours';

class GraphControls extends React.Component {
    constructor() {
        super()
    }

    componentDidMount() {
        this.props.firstLoad()
    }

    render() {
        let { nodeTypes, q, matches, searchNodes, selectNode, clickedNode, view, uiChangeView } = this.props;
        
        return <div className="infoView">
            <div>
  
                <div className='search'>
                    <input type='text' className="form-control" placeholder="Search types, files" onChange={(ev) => searchNodes(ev.target.value)} value={q}/>
                </div>

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


const NodeMatch = ({ onClick, label, variant }) => {
    return <div onClick={onClick}>
        {label}
        <span className="badge badge-light" style={{
            backgroundColor: nodeColor(variant),
            float: 'right'
        }}>{getVariantName(variant)}</span>
    </div>
}


const mapStateToProps = state => {
    return {
        ...state.graph.search,
        nodeTypes: state.graph.nodeTypes,
        clickedNode: state.graph.clickedNode ? state.graph.nodeLookup[state.graph.clickedNode] : null,
        view: state.graph.uiView
    }
}

const mapDispatchToProps = dispatch => {
    return {
        firstLoad:   () => {
            // dispatch(loadInitialFileForTesting())
        },
        searchNodes: (q) => dispatch(searchNodes(q)),
        selectNode:  (id) => dispatch(selectNodeFromSearch(id)),
    }
}

const GraphControlsView = connect(
  mapStateToProps,
  mapDispatchToProps
)(GraphControls)
 
export default GraphControlsView;