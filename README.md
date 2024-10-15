<a name="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="https://github.com/suny-am/bitbucket-cli">
    <img src=".docs/images/logo.png" alt="Logo" width="80" height="80">
  </a>

<h1 align="center">BB</h1>

  <p align="center">
    CLI For Bitbucket Cloud
    <br />
    <!-- <a href="https://github.com/suny-am/bitbucket-cli"><strong>Explore the docs »</strong></a> -->
    <!-- <br /> -->
    <br />
    <a href="https://github.com/suny-am/bitbucket-cli/issues/new?labels=bug&template=bug-report---.md">Report Bug</a>
    ·
    <a href="https://github.com/suny-am/bitbucket-cli/issues/new?labels=enhancement&template=feature-request---.md">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
        <li><a href="#status">Status</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li>
      <a href="#Development">Development</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
    <!-- <li><a href="#acknowledgments">Acknowledgments</a></li> -->
    <li><a href="#references">References</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->

A Bitbucket Cloud compatible CLI that let's you interact with your Bitbucket resources without leaving your terminal.

<p align="right"><a href="#readme-top">🔝</a></p>

### Built With

[![go][go]][go-url]

<p align="right"><a href="#readme-top">🔝</a></p>

### Status

[![FOSSA Status](https://app.fossa.com/api/projects/custom%2B45338%2Fgithub.com%2Fsuny-am%2Fbitbucket-cli.svg?type=shield&issueType=license)](https://app.fossa.com/projects/custom%2B45338%2Fgithub.com%2Fsuny-am%2Fbitbucket-cli?ref=badge_shield&issueType=license)

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- GETTING STARTED -->
## Development

This is an example of how you may give instructions on setting up your project locally.
To get a local copy up and running follow these simple example steps.

### Prerequisites

Testing and developing the application locally requires **go1.22.7** to be available (the latest release to be tested for now).
<br>See the [official documentation](https://go.dev/doc/install) for instructions on installing Go on your platform of choice.

#### macOS Example

  ```sh
  brew install go
  go version
  # example output: 
  # go version go1.22.7 linux/amd64  
  ```

### Installation

#### 1. Clone the repo

   ```sh
   git clone https://github.com/suny-am/bitbucket-cli.git
   ```

#### 2. Build

##### 2.b Local

```

```

##### 2.a Docker

   ```sh
   docker build . -t "bitbucket-cli"
   ```

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- USAGE EXAMPLES -->
## Usage

Currently, only **READ** actions (such as viewing repositories, pull requests and searching for code) is supported, but the plan is to integrate all actions supported by the official Bitbucket REST API.

### Usage examples

#### Repositories

```sh
# List repositories
bb repo list

# List repositories in a workspace
bb repo list -w my-workspace

# List repsitories with a custom limit
bb repo list -l 500

# View a repository
bb repo view -w my-workspace my-repo
```

#### Pullrequests

```sh
# List pullrequests for a repository
bb pr list -w my-workspace -r my-repo

# View a specific pullrequest
bb pr view -w my-workspace -r my-repo my-pullrequest
```

#### Code Search

```sh
# Search for code in a workspace
bb code search -w my-workspace variableName

# Search for code in a repository
bb code search -w my-workspace -r my-repo variableName

# Multiple terms are supported
bb code search -w my-workspace -r my-repo "variable1 variable2 const1"

# As are non ASCII characters
bb code search -w my-workspace -r my-repo "myfunc() => { x % 5 == 0 }"
```

_For more examples, please refer to the [Documentation](https://example.com)_

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- ROADMAP -->
## Feature Roadmap

- [ ] Core commands
  - [x] Root
  - [ ] Workspace
  - [ ] Repository
    - [ ] Branch restrictions
    - [ ] Branching model
    - [x] Commits
    - [ ] Deploy Keys
    - [ ] Downloads
    - [ ] Environments
    - [x] Repositories
    - [x] Pullrequests
    - [ ] Refs
    - [ ] Reports
    - [ ] Source
  - [ ] Project
    - [ ] Branch restrictions
    - [ ] Branching model
    - [ ] Deploy Keys
    - [ ] Issue tracker
    - [ ] Pipelines
    - [ ] Projects
  - [ ] User
    - [x] Pullrequests
  - [ ] Snippets
- [ ] Flags
- [ ] Help Topics
- [ ] Config
- [ ] Auth
- [ ] TBD

See the [open issues](https://github.com/suny-am/bitbucket-cli/issues) for a full list of proposed features (and known issues).

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

### 1. Fork the Project

```sh
gh repo fork suny-am/bitbucket-cli --clone
cd bitbucket-cli
```

### 2. Create your Feature Branch

```sh
git checkout -b feature/aNewCoolFeature
```

### 3. Commit your Changes

```sh
`git commit -m 'Add a new cool feature'
```

### 4. Push to the Branch

```sh
git push origin feature/aNewCoolFeature
```

### 5. Open a Pull Request

```sh
gh pr create 
```

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- LICENSE -->
## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

<p align="right"><a href="#readme-top">🔝</a></p>

<!-- CONTACT -->
## Contact

Your Name - [@bsky_handle](https://bsky.app/profile/bsky_handle) - <visualarea.1@gmail.com>

Project Link: [https://github.com/suny-am/bb](https://github.com/suny-am/bb)

<p align="right"><a href="#readme-top">🔝</a></p>

## References

- [Atlassian Bitbucket Cloud REST API Documentation](https://developer.atlassian.com/cloud/bitbucket/rest/)

<p align="right"><a href="#readme-top">🔝</a></p>
  
<br>
<p align="right"><a href="#readme-top">🔝</a></p>

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/suny-am/bb.svg?style=for-the-badge
[contributors-url]: https://github.com/suny-am/bb/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/suny-am/bb?style=for-the-badge
[forks-url]: https://github.com/suny-am/bb/network/members
[stars-shield]: https://img.shields.io/github/stars/suny-am/bb.svg?style=for-the-badge
[stars-url]: https://github.com/suny-am/bb/stargazers
[issues-shield]: https://img.shields.io/github/issues/suny-am/bb.svg?style=for-the-badge
[issues-url]: https://github.com/suny-am/bb/issues
[license-shield]: https://img.shields.io/github/license/suny-am/bb.svg?style=for-the-badge
[license-url]: https://github.com/suny-am/bb/blob/master/LICENSE
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/carl-sandberg-01070a2b6/
[go]: https://img.shields.io/badge/go-%2300ADD8?style=for-the-badge&logo=go&logoColor=white&logoSize=auto
[go-url]: https://go.dev/
