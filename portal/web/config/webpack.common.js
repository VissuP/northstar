"use strict";
var webpack = require('webpack');
var HtmlWebpackPlugin = require('html-webpack-plugin');
var ExtractTextPlugin = require('extract-text-webpack-plugin');
var helpers = require('./helpers');

module.exports = {
    entry: {
        'polyfills': './polyfills.ts',
        'vendor': './vendor.ts',
        'main': './main.ts'
    },
    resolve: {
        extensions: ['.ts', '.js']
    },

    node: {
        fs: 'empty'
    },

    module: {
        rules: [{
            test: /\.ts$/,
            enforce: 'pre',
            loader: 'tslint-loader',
            options: {

                // can specify a custom config file relative to current directory or with absolute path 
                // 'tslint-custom.json' 
                configFile: './tslint.json',

                // tslint errors are displayed by default as warnings 
                // set emitErrors to true to display them as errors 
                emitErrors: false,

                // tslint does not interrupt the compilation by default 
                // if you want any file with tslint errors to fail 
                // set failOnHint to true 
                failOnHint: false,

                // enables type checked rules like 'for-in-array' 
                // uses tsconfig.json from current working directory 
                typeCheck: false,

                // automatically fix linting errors 
                fix: false,

                // can specify a custom tsconfig file relative to current directory or with absolute path 
                // to be used with type checked rules 
                tsConfigFile: 'tsconfig.json',

                // These options are useful if you want to save output to files 
                // for your continuous integration server 
                fileOutput: {
                    // The directory where each file's report is saved 
                    dir: './tslint-result/',

                    // The extension to use for each report's filename. Defaults to 'txt' 
                    ext: 'xml',

                    // If true, all files are removed from the report directory at the beginning of run 
                    clean: true,

                    // A string to include at the top of every report file. 
                    // Useful for some report formats. 
                    header: '<?xml version="1.0" encoding="utf-8"?>\n<checkstyle version="5.7">',

                    // A string to include at the bottom of every report file. 
                    // Useful for some report formats. 
                    footer: '</checkstyle>'
                }
            }
        }, 
        {
            test: /\.js$/,
            include: helpers.root('node_modules/ngx-vz-cell'),
            use: {
                loader: 'babel-loader',
                options: {
                    presets: ['env']
                }
            }
        },
        {
            test: /\.ts$/,
            loaders: [{
                loader: 'awesome-typescript-loader',
                options: { configFileName: helpers.root('', 'tsconfig.json') }
            },
            { loader: 'angular2-template-loader' },
            {
                loader: 'string-replace-loader',
                query: { search: 'moduleId: module.id,', replace: '' }
            }
            ]
        },
        {
            test: /\.html$/,
            loader: 'html-loader'
        },
        {
            test: /\.(png|jpe?g|gif|svg|woff|woff2|ttf|eot|ico)$/,
            exclude: /node_modules/,
            loader: 'url-loader?limit=1000&name=assets/[name].[hash].[ext]'
        },
        {
            test: /\.css$/,
            use: ['to-string-loader', 'css-loader']
        }
        ]
    },
    plugins: [
        // Workaround for angular/angular#11580
        new webpack.ContextReplacementPlugin(
            // The (\\|\/) piece accounts for path separators in *nix and Windows
            /angular(\\|\/)core(\\|\/)(esm(\\|\/)src|src)(\\|\/)linker/,
            helpers.root('..'), // location of your src
            {} // a map of your routes
        ),
        new HtmlWebpackPlugin({
            template: 'index.html'
        }),
        new webpack.ProvidePlugin({
            nv: 'nvd3'
        })
    ]
}