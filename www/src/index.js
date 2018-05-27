import configureStore from './store';
import ReactDOM from 'react-dom'
import React from 'react'
// import { AppContainer } from 'react-hot-loader';

const rootEl = document.getElementById('root');
const store = configureStore();

// const render = Component => {
//   ReactDOM.render(
//     <AppContainer>
//       <Component store={store} history={history} />
//     </AppContainer>,
//     rootEl,
//   )
// }

// render(App)

import createHistory from 'history/createBrowserHistory'
const history = createHistory()

// Create a reusable render method that we can call more than once
let render = () => {
  // Dynamically import our main App component, and render it
  const App = require("./App").default;
  
  ReactDOM.unmountComponentAtNode(rootEl)
  ReactDOM.render(
      // <Provider store={store}>
      //     <App />
      // </Provider>,
      <App store={store} history={history}/>,
      rootEl
  );
};

if(module.hot) {
  // Support hot reloading of components
  // and display an overlay for runtime errors
  const renderApp = render;
  const renderError = (error) => {
      const RedBox = require("redbox-react").default;
      ReactDOM.render(
          <RedBox error={error} />, 
          rootEl,
      );
  };

  // In development, we wrap the rendering function to catch errors,
  // and if something breaks, log the error and render it to the screen
  render = () => {
      try {
          renderApp();
      }
      catch(error) {
          console.error(error);
          renderError(error);
      }
  };

  // Whenever the App component file or one of its dependencies
  // is changed, re-import the updated component and re-render it
  module.hot.accept("./App", () => {
      setTimeout(render);
  });
}

render();