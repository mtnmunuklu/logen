<h1 align="center">Logen</h1>

<p align="center">
  <a href="https://pkg.go.dev/github.com/mtnmunuklu/logen">
    <img src="https://img.shields.io/badge/%F0%9F%93%9A%20godoc-pkg-informational.svg" alt="Go Doc">
  </a>
  <a href="https://goreportcard.com/report/github.com/mtnmunuklu/logen">
    <img src="https://img.shields.io/badge/%F0%9F%93%9D%20goreport-A+-success.svg" alt="Go Report">
  </a>
  <!-- Other links and badges -->
</p>


Logen is a tool that generates synthetic logs for testing Sigma rules. It reads Sigma rules from files or directories, parses them, and generates synthetic log examples in the "evtx" format using ChatGPT.

## Table of Contents

- [Purpose](#purpose)
- [Installation](#installation)
  - [Normal Installation](#normal-installation)
  - [Docker Installation](#docker-installation)
- [Usage](#usage)
  - [Command-line Flags](#command-line-flags)
  - [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Purpose

Sigma is an open-source project providing a rule format and tools for creating and sharing detection rules for security operations. Logen assists security teams by generating synthetic logs based on Sigma rules, allowing them to test and verify the effectiveness of their rules.

## Installation

To use Logen, you can choose between two installation options:

### Normal Installation

Logen provides precompiled ZIP files for different platforms. Download the appropriate ZIP file for your platform from the following links:

- [Windows](https://github.com/mtnmunuklu/logen/releases/latest/download/logen-windows-latest.zip)
- [Linux](https://github.com/mtnmunuklu/logen/releases/latest/download/logen-ubuntu-latest.zip)
- [macOS](https://github.com/mtnmunuklu/logen/releases/latest/download/logen-macos-latest.zip)

Once downloaded, extract the ZIP file to a directory of your choice. Among the extracted files, you will find the Logen executable.

Ensure that the directory containing the Logen executable is added to your system's PATH environment variable, enabling you to run Logen from any location in the command line.

Note: Logen requires Go to be installed on your system. Download and install Go from the official website: [https://golang.org/dl/](https://golang.org/dl/)

### Docker Installation

Alternatively, you can use Docker to run Logen in a containerized environment. Docker provides a convenient and consistent way to set up and use Logen without worrying about dependencies or system configurations.

To install and set up Logen using Docker, make sure Docker is installed on your system. If not, download and install Docker from the official website: [https://www.docker.com/get-started](https://www.docker.com/get-started)

Once Docker is installed, follow these steps:

1. **Clone the Repository**: If not done already, clone the Logen repository to your local machine:

   ```shell
   git clone https://github.com/mtnmunuklu/logen.git
   ```

2. **Navigate to Docker Directory**: Go to the docker directory inside the cloned repository:
   
   ```shell
   cd tools/docker
   ```

3. **Build Docker Image and Start Container**: Use the setup script to build the Docker image named logen-image:

   ```shell
   go run setup_docker_logen.go -rules <rulesDirectory> -config <configFile> -output <outputDirectory>
   ```

    This script handles building the Docker image and starting the container for you.

That's it! You have successfully installed Logen on your system. Proceed to the  [Usage](#usage) to learn how to use Logen.

If you prefer to build Logen from source, refer to the [Build Instructions](BUILD.md) for detailed steps on building and installing it on your platform.

## Usage

### Command-line Flags

Logen provides several command-line flags for configuring its behavior:

- `filepath`: Name or path of the file or directory to read.
- `config`: Path to the configuration file.
- `filecontent`: Base64-encoded content of the file or directory to read.
- `configcontent`: Base64-encoded content of the configuration file.
- `output`: Output directory for writing files.
- `cs`: Case-sensitive mode.
- `apikey`: API key for ChatGPT.

For more details on available flags, you can use the `-help` flag:
   ```shell
   logen -help
   ```

### Examples

Here are a few examples of using Logen:

- To generate synthetic logs from a Sigma rule file and a configuration file:

   ```shell
   logen -filepath /path/to/sigma/rule.yml -config /path/to/config.yml -apikey your_api_key
   ```
   or
   ```shell
   docker exec logen ./logen -filepath /path/to/sigma/rule.yml -config /path/to/config.yml -apikey your_api_key
   ```

- To generate synthetic logs from Sigma rule content and configuration content:

   ```shell
   logen -filecontent base64_encoded_rule_content -configcontent base64_encoded_config_content -apikey your_api_key
   ```
   or
   ```shell
   docker exec logen ./logen -filecontent base64_encoded_rule_content -configcontent base64_encoded_config_content -apikey your_api_key
   ```

## Contributing

Contributions to Logen are welcome and encouraged! Please read the [contribution guidelines](CONTRIBUTING.md) before making any contributions to the project.

## License

Logen is licensed under the MIT License. See [LICENSE](LICENSE) for the full text of the license.