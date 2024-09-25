![image](https://github.com/user-attachments/assets/c84245e5-5f65-4f7a-8903-b56273867362)

<div style="text-align: center;">
  <h1>
    Introducing the Ultimate Golang Mock Server
  </h1>
</div>

Mockzz-CLI is a command-line tool that enables users to easily create and manage mock APIs along with their associated responses.

### Why Would I use this?

- Easy to set up and install
- Leightweight, interactiv CLI and fast

## Table of Contents

- [Usage Example](#usage-example)
- [Install](#install)
- [License](#license)

## Usage Example
#### Add new mock APIs super easy
![image](./public/create-api.gif)

#### Manage your APIs with ease
![image](./public/manage-apis.gif)

#### Run your mock server with a single command
![image](./public/start-server.gif)

## Install
### Build binary from Source
Install Go: Make sure you have Go installed (version 1.23.0 or later).
Clone the Repository: Clone the source code repository from GitHub.
```bash
git clone https://github.com/CodePawfect/mockzz-cli.git
cd mockzz-cli
```

Build the Binary: Use the Go build tool to create the binary:
```bash
go build -o mockzz
```

Move Binary to Your PATH: Move the generated mockzz binary to a directory in your PATH, so it can be used globally.
```bash
sudo mv mockzz /usr/local/bin/
```

## License
Licensed under [MIT License](./LICENSE.txt)
