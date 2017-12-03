"use strict";
var webpack = require('webpack');
var webpackMerge = require('webpack-merge');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var helpers = require('./helpers');
var CopyWebpackPlugin = require('copy-webpack-plugin');
var commonConfig = require('./webpack.common.js');

const ENV = process.env.NODE_ENV = process.env.ENV = 'production';

module.exports = webpackMerge(commonConfig, {
    devtool: "source-map",

    output: {
        path: helpers.root('dist'),
        publicPath: '/northstar/dist/',
        filename: '[name].[hash].js',
        chunkFilename: '[id].[hash].chunk.js'
    },

    plugins: [
        new webpack.NoEmitOnErrorsPlugin(),
        new webpack.optimize.CommonsChunkPlugin({
            name: ['main', 'vendor', 'polyfills']
        }),
        new CopyWebpackPlugin([
            { from: './assets/img/', to: 'assets/img/' },
            { from: './assets/css/', to: 'assets/css/' },
            { from: './assets/js/', to: 'assets/js/' }
        ]),
        new ExtractTextPlugin({
            filename: "[name].[hash].css"
        }),
        new webpack.DefinePlugin({
            'process.env': {
                'ENV': JSON.stringify(ENV)
            }
        }),
        new webpack.LoaderOptionsPlugin({
            options: {
                htmlLoader: {
                    minimize: false
                }
            }
        })
    ]
});