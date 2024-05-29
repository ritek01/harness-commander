
<h1 align="center">
  <br>
  <img width="1151" alt="image" src="https://github.com/ritek01/harness-cli/assets/139952065/012f928e-b90a-4aa6-8212-4ebad3eb0749">
  <br>
</h1>
<h3 align="center">The proposed CLI tool for Harness consists of several components working together to provide a seamless onboarding and deployment experience.</h3>
<br>


## Key Features

* CLI Interface - Onboard, Sail and Shipped !!
    - Seamlessly Access all quick action properties of Harness with few steps than ever 🐱‍💻
* Authentication Module
    - Authenticate into your Harness account anywhere in the CLI tool 🌎
* Project Initialization Module
    - Initialization of project happens automatically, just say the name you want and done 👍
* Framework Detection Module
    - Automatically detects your deployment type for ready to ship
    - Don't worry you still have options to choose whatever is in your mind 🧑‍💻
* Cross platform
    - MacOS and Linux ready. 🖥️

## How To Use

### Installation
> To clone and run this application, you'll need [Git](https://git-scm.com) and [Go Lang](https://go.dev/doc/install) installed on your computer. From your command line:

1. Download the latest release from the [GitHub releases page](https://github.com/ritek01/harness-cli/releases).

2. Extract the downloaded file to a directory of your choice. It is recommended that you move the extracted file to a folder specified in your system's path for ease of use.

3. If you are using macOS, you can move the harness-upgrade file to the `/usr/local/bin/` directory by running the following command:

```shell
mv harness /usr/local/bin/
```

4. Run the `harness` command to verify that the installation was successful.

### Usage 

```bash
# Login into your account using API key
$ harness login

# Initialize harness in your project 
$ harness init

# Deploy the project
$ harness deploy

# deploy command is a crucial step as this steps does the following :
# - Requests to create a connector or use an existing connector
# - Automatically creates a pipeline
# - Deploys your project into requested server / default Harness server
```

## Tech Stack

[<img src="https://go.dev/images/favicon-gopher.png" style="width: 50px; height: 50px; margin-left: 40px;">](https://go.dev/doc/)[<img src="https://djeqr6to3dedg.cloudfront.net/repo-logos/library/docker/live/logo.png" style="width: 50px; height: 50px;margin-left: 40px;">](https://hub.docker.com/) 

## Credits

This project is created during the Hackweek'24 at Harness Org

- [Ritek Rounak](https://github.com/ritek01)
- [Prakhar Martand](https://github.com/PrakharMartand)

<a href="https://github.com/ritek01/harness-cli/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ritek01/harness-cli" />
</a>
