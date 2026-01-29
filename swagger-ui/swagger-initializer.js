window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "./apispec.json",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    // plugins: [
    //   SwaggerUIBundle.plugins.DownloadUrl
    // ],
    layout: "BaseLayout" // Changed from "StandaloneLayout" to "BaseLayout"
  });

  //</editor-fold>
};
