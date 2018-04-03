import React from 'react';
import ReactDOM from 'react-dom';

import GraphUI from './graph/ui';

document.addEventListener('DOMContentLoaded', function() {
  ReactDOM.render(
    <GraphUI/>,
    document.getElementById('react-mount')
  );
});
