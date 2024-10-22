# Generic Integrator Platform

## Overview

The **Generic Integrator Platform** is a powerful and flexible integration framework that enables seamless connectivity with multiple external service providers. Designed for businesses seeking to streamline their payment processing and third-party service integrations, this platform eliminates the need for cumbersome code changes for each new integration. By leveraging simple configuration files, organizations can efficiently manage interactions with various payment gateways and APIs, enhancing operational efficiency and reducing development time.


## Features

- **Multi-Provider Integration**: Effortlessly connect with a wide range of payment gateways and external services using a single configuration file, simplifying the integration process.
- **Dynamic Configurable Integration**: Manage all integration details through easy-to-update configuration files, minimizing the need for extensive code modifications and allowing for rapid adjustments.
- **Scalable and Extensible Architecture**: Built to accommodate the evolving needs of businesses, the platform supports the addition of new payment methods and integrations without requiring developer intervention, promoting scalability.
- **Secure Data Management**: Ensure the secure handling of sensitive information such as API tokens and authentication credentials, with configurations designed to protect your data.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/ManuelLecaro/generic-integration-platform.git
   cd generic-integration-platform
   ```

2. Install docker and docker-compose
3. Setup the configuration files for your integrations, check payments.toml as the example file to set your integrations configuration
4. Setup the configuration files for the application follow config.toml file
5. Execute 
    ```bash
   make dev
   ``` 
   this will setup the app and its dependencies
7. Go to the api [location](http://localhost:8080/swagger/index.html#/)
8. You can go to [storeDB](http://localhost:2113/web/index.html#/streams)
8. Go and check [mongoDB](http://localhost:8081/)

## Configuration

The platform uses a configuration file (payments.toml) to define how to interact with various payment providers. The file contains details such as base URLs, authentication tokens, endpoints, HTTP methods, and parameter mappings. Here's an example of how to configure a provider:

```toml

[[integrations]]
name = "example_service"
type = "rest"  # Integration type: can be "rest", "grpc", etc.
base_url = "https://api.example.com/v1"
auth_type = "token"  # Possible types: "basic", "token", "oauth"
auth_header = "Authorization"
auth_token = "Bearer TEST_TOKEN"  # Authentication token (can be a template or dynamic value)

# OAuth authentication configuration if needed
[integrations.oauth]
client_id = "your_client_id"
client_secret = "your_client_secret"
token_url = "https://auth.example.com/oauth/token"

# Definition of the provider's available endpoints
[[integrations.endpoints]]
action = "authorize"
method = "POST"
path = "/payment_intents"
description = "Authorize a payment"
# Mapping of the request body parameters (dynamic, based on input variables)
[integrations.endpoints.params]
amount = "{{input.amount}}"
currency = "{{input.currency}}"
payment_method = "{{input.card_number}}"

# Definition of headers specific to this endpoint
[integrations.endpoints.headers]
Authorization = "{{auth_token}}"
Content-Type = "application/json"

# Mapping the response received after processing
[integrations.endpoints.response_mappings]
transaction_id = "{{response.transaction_id}}"
status = "{{response.status}}"

```


## Architecture

We are using CQRS for the payments logic flow, in order to handle the supposed huge amount of payment/refund requests at the same time as we have a need for
almost real time reads for furhter necessities like dashboards and metrics processing. 

We are using EventStoreDB as the event store, so we are ready to get the events
from a payment transaction and also are able to plug more functionalities into the application without scaling issues.

![architecture](img/architecture26.png)


## Pros and Cons of Generic vs. Custom Integration Approaches

| **Aspect**               | **Generic Integration (Config-Based)**                           | **Custom Integration (Code-Based)**                                |
|--------------------------|-----------------------------------------------------------------|-------------------------------------------------------------------|
| **Flexibility**           | **Pro:** High flexibility; add providers via configs without code changes. | **Con:** Less flexible; requires new code for each provider.      |
| **Maintenance**           | **Con:** Configs can become complex.                            | **Pro:** Easier to maintain, each integration is isolated.         |
| **Scalability**           | **Pro:** Scales well; add new providers via configuration.       | **Con:** Requires additional code for each new provider.           |
| **Performance**           | **Con:** Possible overhead from dynamic configs.                | **Pro:** Optimized performance for specific providers.             |
| **Security**              | **Con:** Managing sensitive data in configs can be risky.       | **Pro:** Better security with custom code.                         |
| **Updates & Changes**     | **Pro:** API changes handled with config updates.               | **Con:** Requires code updates for API changes.                    |
| **Development Cost**      | **Con:** Higher initial cost for generic system design.         | **Pro:** Lower initial cost for fewer providers.                   |
| **Extensibility**         | **Pro:** Easily extendable by adding configs.                   | **Con:** Requires new code for each additional provider.           |


# Future work

1. Add terraform configuration to set up necessary components: eks pods for replica and scalling the integration of the API, Route53 for the domain, iam roles, permission groups, ec2 to set up an eventStore instance, dynamoDB as the write model

![architecture](img/AWSARCH.png)

2. Add integration testing and more unit testing for coverage
3. Setup pipelines for CI/CD
4. Set up linting
5. Set up a message bus to integrate with systems that does not need synchronous communication.

To further enhance the **Generic Integrator Platform**, the following features and improvements are planned:

- [ ] **Add Support for gRPC**: Implement gRPC support for improved performance and flexibility in communication between services.
- [ ] **Integrate SOAP Protocol**: Provide support for the SOAP protocol to connect with legacy systems and services.
- [ ] **Implement GraphQL Support**: Enable GraphQL integration for more efficient data querying and manipulation.
- [ ] **Define Integration Extension Format**: Establish a standardized format for extending integrations, making it easier to add new providers and functionalities.
- [ ] **Develop a User Interface**: Create a user-friendly graphical interface to simplify configuration management and improve the overall user experience.
- [ ] **Enhance Testing Suite**: Add integration testing and increase unit test coverage for better reliability and maintainability.
- [ ] **Establish Linting and Code Quality Checks**: Integrate linting tools to maintain code quality and consistency throughout the codebase.
- [ ] **Implement a Message Bus**: Introduce a message bus system to facilitate asynchronous communication between services and improve scalability.