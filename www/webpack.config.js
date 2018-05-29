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
    globalObject: 'this'
  },

  devtool: 'eval-source-map',
  
  plugins: [
    new webpack.LoaderOptionsPlugin({
      options: {
        context: __dirname,
      }
    }),
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
        use: [ 
          'style-loader', 
          {
            loader: 'css-loader',
            options: {
              importLoader: true,
              modules: true,
              localIdentName: "__[name]__[local]___[hash:base64:5]"
            }
          }
        ]
      },
      { 
        test: /worker\.js$/,
        use: ['babel-loader', 'worker-loader']
      },
      {
        test: /\.js$/,
        exclude: /(node_modules|vendor)/,
        use: ['babel-loader'],
      }
    ],
  },
};