"use strict";
var webpackMerge = require('webpack-merge');
var commonConfig = require('./webpack.common.js');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var helpers = require('./helpers');
var CopyWebpackPlugin = require('copy-webpack-plugin');
var webpack = require('webpack');

module.exports = webpackMerge(commonConfig, {
    devtool: 'source-map',

    output: {
        path: helpers.root('dist'),
        publicPath: '/northstar',
        filename: '[name].[hash].js',
        chunkFilename: '[id].[hash].chunk.js'
    },


    plugins: [
        new webpack.optimize.CommonsChunkPlugin({
            name: ['main', 'vendor', 'polyfills']
        }),
        new ExtractTextPlugin({
            filename: "[name].[hash].css"
        }),
        new CopyWebpackPlugin([
            { from: './assets/img/', to: 'assets/img/' },
            { from: './assets/css/', to: 'assets/css/' },
            { from: './assets/js/', to: 'assets/js/' }
        ]),
    ],

    devServer: {
        historyApiFallback: true,
        contentBase: helpers.root(''),
        watchContentBase: true,
        stats: 'minimal',
        port: 8081,
        proxy: {
            'http://localhost:8081/users/**': {
                target: 'http://localhost:8080'
            },
            'http://localhost:8081/ns/v1/**':{
                target: 'http://localhost:8080'
            },
            'http://localhost:8081/northstar/app/shared/images/**':{
                target: 'http://localhost:8080'
            }
            // 'ws://localhost:8081/ns/v1/connections/':{
            //     target: 'ws://localhost:8080/ns/v1/connections/'
            // }
        }
    }
});