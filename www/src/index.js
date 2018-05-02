import App from './app';
import ReactDOM from 'react-dom'
import React from 'react'
import  { AppContainer } from 'react-hot-loader';


const render = Component => {
  ReactDOM.render(
    <AppContainer>
      <Component />
    </AppContainer>,
    document.getElementById('root'),
  )
}

render(App)

if (module.hot) {
  module.hot.accept();

  module.hot.dispose((data) => {
    data.store = store;
  });
}


// Webpack Hot Module Replacement API
if (module.hot) {
  module.hot.accept('./app', () => {
    // if you are using harmony modules ({modules:false})
    render(App)
    // in all other cases - re-require App manually
    // render(require('./app'))
  })
}