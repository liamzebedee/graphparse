const path = require('path');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const webpack = require('webpack');
const FlowWebpackPlugin = require('flow-webpack-plugin')
const serve = require('./webpack.serve');

module.exports = {
  mode: 'development',
  target: 'web',
  serve,

  context: path.join(__dirname, 'src'),

  entry: './index',

  output: {
    path: path.resolve(__dirname, "dist"),
    filename: 'bundle.js',
    publicPath: '/',
  },

  devtool: 'eval-source-map',
  
  plugins: [
    // new webpack.HotModuleReplacementPlugin(),
    new FlowWebpackPlugin(),
    new HtmlWebpackPlugin({
      title: "Basemap",
      template: './index.ejs',
    }),
    new webpack.DefinePlugin({
      "process.env": {
         NODE_ENV: JSON.stringify("production") 
       }
    }),
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