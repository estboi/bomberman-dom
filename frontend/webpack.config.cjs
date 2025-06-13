const path = require("path");
const HtmlWebpackPlugin = require("html-webpack-plugin");

module.exports = {
  entry: "./index.js", // Entry point of the application
  output: {
    path: path.resolve(__dirname, "dist"), // Output directory path
    filename: "bundle.js", // Output bundle file name
    publicPath: "/",
  },
  module: {
    rules: [
      {
        test: /\.(js|jsx|ts|tsx)$/,
        exclude: /node_modules/, // Exclude node_modules folder from processing
        use: {
          loader: "babel-loader", // Use Babel to transpile JavaScript
          options: {
            presets: [
              [
                "@babel/preset-typescript",
                {
                  jsxPragma: "Gachi.createElement",
                  jsxPragmaFrag: "Gachi.Fragment",
                },
              ],
            ],
            plugins: [
              [
                "@babel/plugin-transform-react-jsx",
                {
                  pragma: "Gachi.createElement",
                  pragmaFrag: "Gachi.Fragment",
                },
              ],
            ],
          },
        },
      },
      {
        test: /\.css$/,
        use: ["style-loader", "css-loader"],
      },
      {
        test: /\.(png|svg|jpeg|jpg|gif)$/,
        use: ["file-loader"],
      },
    ],
  },
  resolve: {
    extensions: [".js", ".jsx", ".ts", ".tsx"], // Add support for resolving .js and .mjs extensions
  },
  devServer: {
    static: path.join(__dirname, "./Public"), // Path to assets
    historyApiFallback: true,
    port: 3000,
  },
  plugins: [
    new HtmlWebpackPlugin({
      template: path.resolve(__dirname, "./Public/index.html"), // Path to your HTML template
      filename: "index.html", // Output HTML file name
    }),
  ],
};
