const webpackDevServer = require('webpack-dev-server');
const webpack = require('webpack');

const config = require('./webpack.config.js');
const options = {
  contentBase: './dist',
  hot: true,
  host: 'localhost'
};

webpackDevServer.addDevServerEntrypoints(config, options);
const compiler = webpack(config);
const server = new webpackDevServer(compiler, options);

server.listen(3000, 'localhost', () => {
  console.log('dev server listening on port 3000');
});

// const express = require('express');
// const webpackDevMiddleware = require('webpack-dev-middleware');
// const webpack = require('webpack');
// const webpackConfig = require('./webpack.config.js');
// const app = express();

// const compiler = webpack(webpackConfig);

// app.use(webpackDevMiddleware(compiler, {
//   hot: true,
//   filename: 'bundle.js',
//   publicPath: '/',
//   stats: {
//     colors: true,
//   },
//   historyApiFallback: true,
// }));

// app.use(express.static(__dirname + '/www'));

// const server = app.listen(3000, function() {
//   const host = server.address().address;
//   const port = server.address().port;
//   console.log('app listening at http://%s:%s', host, port);
// });
