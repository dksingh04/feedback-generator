<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="./public/img/feedback3.jpg" alt="Project logo"></a>
</p>

<h3 align="center">Interview Feedback Generator</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/kylelobo/The-Documentation-Compendium.svg)](https://github.com/kylelobo/The-Documentation-Compendium/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/kylelobo/The-Documentation-Compendium.svg)](https://github.com/kylelobo/The-Documentation-Compendium/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> Project can be used for generating formatted feedback for any technical discussion you do, simply by answering few questions.
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [Authors](#authors)

## 🧐 About <a name = "about"></a>

Project will help you to generate formatted feedback by answering questions, it will save around 10-15 mins of writing feedback of Interviews taken.

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

In order to run the code locally, minimal requirement is to have Go installed, download and follow the instruction for [installation](https://golang.org/doc/install). Using google protobuf for grpc API's which expose Creat, Read, Delete and Generat report stored in MongoDB.

There are two flavor of this project, you can use simple UI and answer few questions which will generate formatted report without interacting to DB, but if you want to use the project fully you can go with MongoDB, which requires [MongoDB](https://docs.mongodb.com/manual/installation/) to be setup locally. You can change DB configuration from resources/config.yaml file.

You can interact to these API, by writing your own client and UI. Use protoc to generate client stub in whatever language you prefer and utilize the API's.

### Installing



## 🔧 Running the tests <a name = "tests"></a>

TODO in progress

### Break down into end to end tests

TODO in progress

```
Give an example
```

## 🎈 Usage

## By default it builds the binary for Mac/Linux

TODO update make to build the binary for other OS.

### Build server
make build

### Build client
make build-client

### Build cli client
make build-cli-client

Above command will create the binary and to run them simply execute the binary as below
### Start server 
./feedback-genreator

### Start client
./feedback-client

### Start CLI client
./feedback-cli-client


## 🚀 Deployment
No deployment Instruction, If you have Binary file and just want to run locally and generate feedback, perform following steps

Open [localhost] (http://localhost:8080) on browser.

If you want to change port, make the changes in resources/config.yaml file

## ⛏️ Built Using

- [MongoDB](https://www.mongodb.com/) - MongoDB Database
- [Protobuf](https://developers.google.com/protocol-buffers/docs/gotutorial) - Google Protobuf
- [GoLang](https://golang.org/) - GoLang
- [HTML, JavaScript and JQuery](https://nodejs.org/en/) - HTML, JavaScript and JQuery

## ✍️ Authors

- [@dksingh04](https://github.com/dksingh04) - Deepak Singh