import mobx, { observable, computed } from 'mobx';
import { observer } from "mobx-react";
import _ from 'underscore';
import Fuse from 'fuse.js';

import graphdataFromFile from '../../graph.json';
var graphdata = graphdataFromFile;

import { renderGraph } from '../graph';

export default class GraphStore {
    @observable hoverId = null;
    @observable highlightedNodes = null;
    @observable selectedNodeIds = [];
    @observable nodeTypesHidden = [];

    @observable updateGraph = new Date;
    nodes = graphdata.nodes;
    edges = graphdata.edges;
  
    @computed get selectionInfo() {
      let key = this.hoverId || this.selectionId;
      return graphdata.nodesLookup[key]
    }
  
    @computed.struct get selectedNodes() {
      return this.selectedNodeIds.map(id => { return graphdata.nodesLookup[id] })
    }
  
    loadGraph() {
      fetch(`http://localhost:8081/graph`)
      .then(response => response.json())
      .then(data => {
        graphdata = data;
        renderGraph(graphdata.nodes, graphdata.edges)
      })
    }
  
    mouseoverNode(id) {
      this.hoverId = id;
    }
  
    mouseoutNode(id) {
      this.hoverId = null;
    }
  
    selectNode(id) {
      if(!_.contains(this.selectedNodeIds, id))
        this.selectedNodeIds.push(id)
    }
  
    clearSelection(id) {
      this.selectionId = null;
    }
  
    searchAndHighlightNodes(query) {
      if(query == "") {
        this.highlightedNodes = null;
        return
      }
      let matches = [];
  
      function splitId(id) {
        return [id]
      }
  
      let searchData = graphdata.nodes.map(node => {
        return {
          id: node.id, 
          label: splitId(node.label),
        }
      })
  
      var options = {
        shouldSort: true,
        findAllMatches: true,
        includeScore: true,
        threshold: 0.3,
        location: 0,
        distance: 100,
        maxPatternLength: 32,
        minMatchCharLength: 2,
        keys: [
          "label"
        ]
      };
  
      let fuse = new Fuse(searchData, options);
      let results = fuse.search(query)
      this.highlightedNodes = results.filter(el => el.score < 0.7)
    }
  
    getNodeHighlight(id) {
      return this.highlightedNodes.filter(el => el.item.id == id)[0];
    }
  
    getCodeThread() {
      let from = this.selectedNodeIds[0]
      let to = this.selectedNodeIds[1]
  
      fetch(`http://localhost:8081/graph/thread/${from}/${to}`)
      .then(response => response.json())
      .then(data => {
        renderGraph(graphdata.nodes, graphdata.edges)        
      })
    }
  
    toggleNodeTypeFilter(nodeTypeIdx) {
      if(_.contains(this.nodeTypesHidden, nodeTypeIdx)) {
        this.nodeTypesHidden = _.filter(this.nodeTypesHidden, (el) => el != nodeTypeIdx)
      } else {
        this.nodeTypesHidden.push(nodeTypeIdx)
      }
      this.getFilteredGraph()
    }

    getFilteredGraph() {
      let q = JSON.stringify({
        nodeTypesHidden: this.nodeTypesHidden
      });

      fetch(`http://localhost:8081/graph/filtered?q=${q}`)
      .then(response => response.json())
      .then(data => {
        this.nodes = data.nodes;
        this.edges = data.edges;
        this.updateGraph = new Date;  
      })
    }
  }