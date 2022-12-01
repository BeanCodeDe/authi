# Authi
A lightweight authentication service written in go

![Build](https://img.shields.io/github/workflow/status/BeanCodeDe/authi/MainPipeline.svg)
![License](https://img.shields.io/github/license/BeanCodeDe/authi.svg)
![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)
![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/BeanCodeDe/authi.svg)

## About
Authi is a lightweight authentication service written in go. It covers the basic use cases of creating and deleting users as well as refresh tokens and update passwords.

## Usage

<html>
    <head>
        <!-- Load the latest Swagger UI code and style from npm using unpkg.com -->
        <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
        <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css"/>
        <title>My New API</title>
    </head>
    <body>
        <div id="swagger-ui"></div> <!-- Div to hold the UI component -->
        <script>
            window.onload = function () {
                // Begin Swagger UI call region
                const ui = SwaggerUIBundle({
                    url: "api/open-api.yaml", //Location of Open API spec in the repo
                    dom_id: '#swagger-ui',
                    deepLinking: true,
                    presets: [
                        SwaggerUIBundle.presets.apis,
                        SwaggerUIBundle.SwaggerUIStandalonePreset
                    ],
                    plugins: [
                        SwaggerUIBundle.plugins.DownloadUrl
                    ],
                })
                window.ui = ui
            }
        </script>
    </body>
