const VueLoaderPlugin = require('vue-loader/lib/plugin')

module.exports = {
  chainWebpack: config => {
    config.module
      .rule('vue')
      .use('vue-loader')
      .loader('vue-loader')
      .options({
        transformAssetUrls: {
          video: ['src', 'poster'],
          source: 'src',
          img: 'src',
          image: 'xlink:href',
          'b-img': 'src',
          'b-img-lazy': ['src', 'blank-src'],
          'b-card': 'img-src',
          'b-card-img': 'src',
          'b-card-img-lazy': ['src', 'blank-src'],
          'b-carousel-slide': 'img-src',
          'b-embed': 'src',
          'Project': 'logo',
          'CareerEntry': 'logo',
        },
      })
  },
  publicPath: '/fancy'
}
// vue.config.js file to be place in the root of your repository
