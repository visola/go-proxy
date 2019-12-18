const CopyPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

module.exports = {
  entry: './src/js/index.js',

  output: {
    filename: 'main.js',
    path: path.resolve(__dirname, 'dist'),
  },

  plugins: [
    new CopyPlugin([
      { from: 'src/css', to: path.resolve(__dirname, 'dist/css') },
    ]),

    new HtmlWebpackPlugin({
      template: 'src/html/index.html',
      title: 'go-proxy admin',
    })
  ]
};
