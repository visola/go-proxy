const CopyPlugin = require('copy-webpack-plugin');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const path = require('path');

const mode = process.env.NODE_ENV || 'development';
const prod = mode === 'production';

module.exports = {
  devtool: prod ? 'none' : 'inline-source-map', 
  entry: './src/js/main.js',

  module: {
    rules: [
      {
        test: /\.svelte$/,
        use: {
          loader: 'svelte-loader',
          options: {
            emitCss: true,
            hotReload: true
          }
        }
      },
    ]
  },

  output: {
    chunkFilename: '[name].[id].js',
    filename: '[name].js',
    path: path.resolve(__dirname, 'dist'),
  },

  plugins: [
    new CopyPlugin([
      { from: 'src/css', to: path.resolve(__dirname, 'dist/css') },
    ]),

    new HtmlWebpackPlugin({
      template: 'src/html/index.html',
      title: 'go-proxy admin',
    }),
  ],

  resolve: {
    alias: {
      svelte: path.resolve('node_modules', 'svelte')
    },
    extensions: ['.mjs', '.js', '.svelte'],
    mainFields: ['svelte', 'browser', 'module', 'main']
  },
};
