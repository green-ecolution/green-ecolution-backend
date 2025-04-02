<p>
  <a href=""><img src="https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white"/></a>
  <a href=""><img src="https://raw.githubusercontent.com/green-ecolution/green-ecolution-backend/badges/.badges/develop/coverage.svg"/></a>
  <a href="https://pkg.go.dev/github.com/green-ecolution/green-ecolution-backend"><img src="https://pkg.go.dev/badge/github.com/green-ecolution/green-ecolution-backend.svg" alt="Go Reference"></a>
</p>

# Green Ecolution Backend

![image](https://github.com/user-attachments/assets/c69e28cd-44ec-44e7-8d64-447bad7f4fd9)

Smart irrigation is needed to save water, staff and costs. This project is the server-side for Green Ecolution. For Frontend please refer to [Green Ecolution Frontend](https://github.com/green-ecolution/green-ecolution-frontend). The backend provides an interface to interact between the website and the database. The backend retains data about

- trees
- tree clusters
- sensors
- watering plans
- vehicles
- users

In the current setup sensors are connected to an ESP32 with an integrated LoRaWAN module.
Sensor data is send using LoraWAN to a MQTT-Gateway and then to the server to further process the data.

While the project is created in collaboration with the local green space management (TBZ Flensburg), this software aims to be applicable to other cities. It was initially developed at the University of Applied Sciences Flensburg as a research project within the Applied Computer Science masters degree program.

For more information please refer to:
- [Project website](https://green-ecolution.de/)
- [University of Applied Sciences Flensburg](https://hs-flensburg.de/en)
- [Try it out | Demo](https://demo.green-ecolution.de)

## Quick start guide

This is a simple way to start the application locally without a lot of setup.
The first command you need to execute is to download the docker compose configuration file from the specified GitHub repository:
`wget https://raw.githubusercontent.com/green-ecolution/green-ecolution-backend/refs/heads/develop/.docker/docker-compose.yaml`

The next command builds and starts the application containers as defined in the configuration file.
`docker compose up`

Now you can run the application with `make run/live`. For further information you can take a look at the setup section down below.

## Project structure

```
.
├── config       <- configuration files
│   ├── app.go
│   └── ...
├── internal
│   ├── entities <- domain entities (models)
│   ├── server   <- server setup (http, grpc, etc.)
│   ├── service  <- business logic (services)
│   └── storage  <- storage logic repository (database, cache, etc.)
└ main.go
```

## Technologies

- [Golang](https://go.dev/) as the main programming language
- [env](https://github.com/caarlos0/env) for environment variables
- [godotenv](https://github.com/joho/godotenv) for loading environment variables from a `.env` file
- [fiber](https://docs.gofiber.io/) for the web framework
- [testify](https://github.com/stretchr/testify) for testing

## Architecture

### Clean Architecture

The project is structured following the principles of the [Clean Architecture]. The main idea is to separate the business logic from the infrastructure and the framework. This way, the business logic is independent of the framework and can be easily tested. The framework is only used to connect the business logic with the outside world. The business logic is divided into three layers: entities, use cases, and interfaces. The entities layer contains the domain models. The use cases layer contains the business logic. The interfaces layer contains the interfaces that the use cases need to interact with the outside world. On top of it the project is structured following the principles of the [Layered Architecture].

Inside the `internal` folder, there are three main packages: `entities`, `service`, and `storage`. The `entities` package contains the domain models. The `service` package contains the business logic. The `storage` package contains the repository logic to interact with the database, cache, etc.

Inside the `internal` folder, there is a `server` package that contains the server setup. The server setup is responsible for setting up the server (http, grpc, etc.) and connecting the business logic with the outside world.

[Clean Architecture]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[Layered Architecture]: https://medium.com/@shershnev/layered-architecture-implementation-in-golang-6318a72c1e10

## Local development

### Requirements

- [Golang](https://go.dev/) as the main programming language
- [Air](https://github.com/air-verse/air) for live reload
- [Mockery](https://github.com/vektra/mockery) for mocking interfaces. Use version `v2.43.2`
- [Make](https://www.gnu.org/software/make/) to execute Makefile
- [docker](https://github.com/docker) for containers
- [docker-compose](https://github.com/docker/compose) to manage containers

### Setup

To download all needed tools use

```bash
make setup
```

Then to generate code use

```bash
make generate
```

To run a local database you can use a preconfigured .yaml file.

```bash
docker compose up -d
```

### Run

To run the project, you need to execute the following command:

**With live reload**

```bash
make run/live
```

**Without live reload**

```bash
make run
```

### Test

Before running the tests, you need to create the mock files. To create the mock files, you need to execute the following command:

```bash
go install github.com/vektra/mockery/v2@v2.43.2 # install Mockery
mockery # create mock files
```

To run the tests, you need to execute the following command:

```bash
go test ./...
```

**NOTE:** Mockery is used to generate mocks for interfaces. The mocks are generated in the `_mocks` folder. To specify the output folder or to add created interfaces to the mocks, you can edit the `mockery.yml` file. The `mockery.yml` file is used to configure the behavior of Mockery. Running `go generate` will execute Mockery and generate the mocks. Also when running Air, the mocks will be generated automatically.

### How to contribute

If you want to contribute to the project please follow this guideline:

- Fork the project.
- Create a topic branch from develop.
- Make some commits to improve the project.
- Push this branch to your GitHub project.
- Open a Pull Request on GitHub.
- Discuss, and optionally continue committing.
- The project owner merges or closes the Pull Request.

Please refer to naming conventions for branches [Medium Article](https://medium.com/@abhay.pixolo/naming-conventions-for-git-branches-a-cheatsheet-8549feca2534).
