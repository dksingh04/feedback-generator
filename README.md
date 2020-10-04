<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="./public/img/feedback4.jpeg" alt="Project logo"></a>
</p>

<h3 align="center">Interview Feedback Generator</h3>

<!--div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div-->

---

<p align="center"> Generating formatted feedback for any technical discussion by simply answering few questions.
    <br> 
</p>

## 🧐 About <a name = "about"></a>

Project will help you to generate formatted feedback by answering few sample questions, it will save around 5-10 mins of your precious time.

## 🏁 Getting Started <a name = "getting_started"></a>

TODO add notes

### Prerequisites

In order to run the code locally, minimal requirement is to have Go (higher than 1.11, and go module for dependencies) installed, download and follow the instruction for [installation](https://golang.org/doc/install). Using google protobuf for grpc API's which expose Creat, Read, Delete and Generat report stored in MongoDB.

There are two flavor of this project, you can use simple UI and answer few questions which will generate formatted report without interacting to DB, but if you want to use the project fully you can go with MongoDB, which requires [MongoDB](https://docs.mongodb.com/manual/installation/) to be setup locally or you can change DB configuration from resources/config.yaml file wherever your mongodb is setup.

You can interact to these API, by writing your own client and UI. Use **protoc** to generate client stub in whatever language you prefer and utilize the API's.


## 🎈 Usage

## By default it builds the binary for Mac/Linux

TODO update to add entry in make file to build binaries for other OS.

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
