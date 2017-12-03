"use strict";
var webpack = require('webpack');
var webpackMerge = require('webpack-merge');
var helpers = require('./helpers');
var commonConfig = require('./webpack.common.js');
var DefinePlugin = require('webpack/lib/DefinePlugin');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var LoaderOptionsPlugin = require('webpack/lib/LoaderOptionsPlugin');

var ENV = process.env.NODE_ENV = process.env.ENV = 'test';

module.exports = webpackMerge(commonConfig, {
    devtool: 'inline-source-map',

    module: {
        rules: [{
                enforce: 'pre',
                test: /\.js$/,
                loader: 'source-map-loader',
                exclude: [
                    // these packages have problems with their sourcemaps
                    helpers.root('node_modules/rxjs'),
                    helpers.root('node_modules/@angular')
                ]
            },
            {
                enforce: 'post',
                test: /\.(js|ts)$/,
                loader: 'istanbul-instrumenter-loader',
                include: helpers.root('app'),
                exclude: [
                    /\.(e2e|spec)\.ts$/,
                    /node_modules/
                ]
            }
        ]
    },

    /**
     * Add additional plugins to the compiler.
     *
     * See: http://webpack.github.io/docs/configuration.html#plugins
     */
    plugins: [
        /**
         * Plugin: DefinePlugin
         * Description: Define free variables.
         * Useful for having development builds with debug logging or adding global constants.
         *
         * Environment helpers
         *
         * See: https://webpack.github.io/docs/list-of-plugins.html#defineplugin
         */
        // NOTE: when adding more properties make sure you include them in custom-typings.d.ts
        new DefinePlugin({
            'ENV': JSON.stringify(ENV),
            'HMR': false,
            'process.env': {
                'ENV': JSON.stringify(ENV),
                'NODE_ENV': JSON.stringify(ENV),
                'HMR': false,
            }
        }),
        /**
         * Plugin LoaderOptionsPlugin (experimental)
         *
         * See: https://gist.github.com/sokra/27b24881210b56bbaff7
         */
        new LoaderOptionsPlugin({
            debug: false,
            options: {
                // legacy options go here
            }
        }),

    ],

    /**
     * Disable performance hints
     *
     * See: https://github.com/a-tarasyuk/rr-boilerplate/blob/master/webpack/dev.config.babel.js#L41
     */
    performance: {
        hints: false
    },
    /**
     * Include polyfills or mocks for various node stuff
     * Description: Node configuration
     *
     * See: https://webpack.github.io/docs/configuration.html#node
     */
    node: {
        global: true,
        process: false,
        crypto: 'empty',
        module: false,
        clearImmediate: false,
        setImmediate: false
    }

});