<div align="center">
  <a href="https://go.dev/">
    <img src="https://img.shields.io/badge/Go-v1.21-brightgreen.svg" alt="go version">
  </a>
</div>

<div align="center">

  <h1>Serpent</h1>
  <h3>Control snake. Eat food. Create chaos.</h3>

</div>

## Contents
- [**What is Serpent?**](#-what-is-serpent-)
- **[Installation](#installation)**
  - [Build from source](#build-from-source-)
  - [Precompiled binaries](#precompiled-binaries-)
- [**Usage**](#usage)
  - [Starting the game](#starting-the-game-)
  - [Playing Serpent](#playing-serpent-)
  - [Kubernetes interaction](#kubernetes-interaction-)
- [**Contribute**](#contribute-)
- [**Acknowledgements**](#acknowledgments)

## ⭐ What is Serpent? ⭐

Serpent lets you play snake while wrecking havock in your Kubernetes cluster. Have fun while you can.

### How does it work? 🤔

Each piece of food you eat corresponds to a pod in your cluster (I left out kube-system though..).

## Installation

### Homebrew 🍺

To install Serpent using Homebrew, you can run the following commands:

```sh
brew tap deggja/serpent https://github.com/deggja/serpent
brew install serpent
```

### Build from source 💻

To build Serpent from the source, you need a working Go environment with version 1.21 or higher. Follow these steps:

```sh
git clone https://github.com/deggja/serpent.git --depth 1
cd serpent
go build -o serpent
```

## Usage

### Starting the game

To start the game, simply run the compiled binary:

```sh
./serpent
```

This will run the game in `default mode`. The snake will only eat resources of type `pod` and avoid system critical workloads in `kube-system`.

To specify a configuration file, use the `--config/-c` flag.

```sh
./serpent --config config.json
```
This will run the game in `config mode`. The snake will eat all resource types in all namespaces defined in the configuration file.

### Example Configuration File

```json
{
    "resource_types": ["pods", "replicasets", "deployments", "services"],
    "namespaces": {
        "include": ["grafana", "default", "netfetch", "podinfo", "workloads"],
        "exclude": ["kube-system"]
    }
```

## Playing Serpent

Use the arrow keys to navigate the snake around the screen:

| Key | Action               |
|-----------------|----------------------|
| Arrow up        | Move up              |
| Arrow down      | Move down            |
| Arrow left      | Move left            |
| Arrow right     | Move right           |
| Space           | Pause or Resume      |
| CTRL + C        | Quit the game        |

[![asciicast](https://asciinema.org/a/Q4usmR4HB8LhHojJA9qJeQmdX.svg)](https://asciinema.org/a/Q4usmR4HB8LhHojJA9qJeQmdX)

## Kubernetes interaction

Serpent will needs access to a Kubernetes cluster. Ensure your `kubeconfig` is set up correctly before starting the game. The application currently expects the default kubeconfig or a kubeconfig environment variable.

As you play and the pods are deleted, Serpent will log its actions to a `chaos.log` file for your review.

## Contribute 🔨

Feel free to dive in! [Open an issue](https://github.com/deggja/serpent/issues) or submit PRs.

## Acknowledgments

This project utilizes [Termloop](https://github.com/JoelOtter/termloop), a simple Go library for creating terminal-based games. Thanks to the creators and contributors of Termloop for providing such a versatile tool.

## License

Serpent is released under the MIT License. Check out the [LICENSE](https://github.com/deggja/serpent/LICENSE) file for more information.
