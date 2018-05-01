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
  devtool: "inline-source-map",
  
  plugins: [
    new HtmlWebpackPlugin({
      template: './index.html',
      title: "Basemap"
    }),
    new webpack.HotModuleReplacementPlugin(),
    new webpack.SourceMapDevToolPlugin({
      filename: "[file].map"
    }),
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
  // content: [path.join(__dirname, 'src')],

  // context: path.join(__dirname, 'src'),
  add: (app, middleware, options) => {
    const historyOptions = {
    };

    app.use(convert(history(historyOptions)));
  }
};