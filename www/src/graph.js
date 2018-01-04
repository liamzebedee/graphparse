import 'script-loader!./vendor/d3.v4.min.js';
import 'script-loader!./vendor/viz-lite.js';
import 'script-loader!./vendor/d3-graphviz.min.js';
import graphdata from 'raw-loader!../graph.dot';

document.addEventListener('DOMContentLoaded', function() {
  let g = d3.select("#graph");
  g.graphviz().renderDot(graphdata).totalMemory(16777216 * 2) // 32mb memory
  g.attr('width', '100%').attr('height', '100%')
  // g.select('.node').
});
