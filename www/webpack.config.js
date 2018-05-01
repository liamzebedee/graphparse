const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');

module.exports = {
  mode: 'development',
  target: 'web',

  context: path.join(__dirname, 'src'),

  // entry: ['./index', 'webpack/hot/only-dev-server', 'webpack-dev-server/client?http://0.0.0.0:3001'],
  entry: './index',

  output: {
    path: path.resolve(__dirname, "dist"),
    filename: 'bundle.js',
    publicPath: '/',
  },

  // devtool: "inline-cheap-module-source-map",
  // devtool: "inline-source-map",
  devtool: 'eval-source-map',
  
  plugins: [
    new HtmlWebpackPlugin({
      title: "Basemap",
      template: './index.ejs',
    }),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.DefinePlugin({
      "process.env": {
         NODE_ENV: JSON.stringify("production") 
       }
    })
  ],

  module: {
    rules: [
      {
        test: /\.css$/,
        use: [ 'style-loader', 'css-loader' ]
      },
      {
        test: /\.js$/,
        exclude: /(node_modules|vendor)/,
        use: ['babel-loader'],
      }
    ],
  },
};


const history = require('connect-history-api-fallback');
const convert = require('koa-connect');

module.exports.serve = {
  add: (app, middleware, options) => {
    const historyOptions = {};
    app.use(convert(history(historyOptions)));
  }
};