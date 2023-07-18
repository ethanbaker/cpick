<!--
  Created by: Ethan Baker (contact@ethanbaker.dev)
  
  Adapted from:
    https://github.com/othneildrew/Best-README-Template/

Here are different preset "variables" that you can search and replace in this template.
-->

<div id="top"></div>


<!-- PROJECT SHIELDS/BUTTONS -->
![1.2.1](https://img.shields.io/badge/status-1.2.1-red)
[![GoDoc](https://godoc.org/github.com/ethanbaker/cpick?status.svg)](https://godoc.org/github.com/ethanbaker/cpick)
[![Go Report Card](https://goreportcard.com/badge/github.com/ethanbaker/cpick)](https://goreportcard.com/report/github.com/ethanbaker/cpick)

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br><br><br>
<div align="center">
  <a href="https://github.com/ethanbaker/cpick">
    <img src="./docs/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Cpick</h3>

  <p align="center">
    An extensive color picker for the terminal
  </p>
</div>


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
    <li><a href="#contact">Contact</a></li>
  </ol>
</details>


<!-- ABOUT -->
## About

![Project demonstration image][product-screenshot]

Cpick is an interactive color picker in the terminal. You can run Cpick in any true color terminal, and you can see thousands of unique colors, either from preset values or gradients. Each color has its own formats in many different forms, including RGB, HSV, CMYK, and more. 

<p align="right">(<a href="#top">back to top</a>)</p>


### Built With

* [Golang](https://go.dev/)
* [Cview](https://code.rocketnine.space/tslocum/cview)
* [Tcell](https://github.com/gdamore/tcell)

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- GETTING STARTED -->
## Getting Started

Cpick depends on [Tcell](github.com/gdamore/tcell) and [Colors](github.com/ethanbaker/colors).

Cpick also uses [Cview](gitlab.com/tslocum/cview). However, Cview currently has a feature where the table cells are drawn with a required space in between cells, which ruins the aesthetic of Cpick. In order to fix this, a forked version of Cview is used within the Cpick package that has the required fixes. This may lead to bugs and issues in older versions of Cview that are not able to be readily fixed until Cview adopts required features.


### Prerequisites

* Golang is installed
* Git is installed
* Your terminal is [true color compatible](https://unix.stackexchange.com/questions/450365/check-if-terminal-supports-24-bit-true-color)


### Installation

#### Golang 

To include the Cpick package in your own Golang application, include the line `import "github.com/ethanbaker/cpick"`.

To make a command based on the package, run `go install` in the `cmd/cpick` directory. Command usage can be found in the [documentation][documentation-url]. If you have difficulties compiling and installing the application, follow [these comprehensive steps](https://go.dev/doc/tutorial/compile-install) in the official Golang documentation.

Cpick also utilizes [cmdtab](https://github.com/rwxrob/cmdtab), which offers tab completion for Cpick. If you wish to enable tab completion, add the line `complete -C cpick cpick` to your `~/.bashrc` file.

#### Arch Linux (Arch User Repository)

Visit the [AUR page](https://aur.archlinux.org) to download Cpick to your arch system. If you don't know how to download a package from the Arch User Repository, you can follow directions from the [Arch wiki](https://wiki.archlinux.org/index.php/Arch_User_Repository#Installing_packages).

#### Ubuntu/Debian

Coming soon!

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- USAGE EXAMPLES -->
## Usage

Cpick can be used to select any color and output the corresponding color format. For example, entering the command `cpick hex` will output a hexadecimal color (ex: `#123abc`) when a color is selected. Many different color formats exist, such as ansi escape sequences, rgb values, or even color names.

For more examples, please refer to the [documentation][documentation-url].

#### Cpick Settings

Cpick comes with two boolean settings: sandbox and testing.

* Sandbox

  Sandbox determines whether or not Cpick will return a color. Sandbox is useful for looking at colors that will be used somewhere else, like for a website. If sandbox is set to `true`, Cpick will not return any value but will instead cycle between different color tables. Setting sandbox to `false` means a color will be returned. 

* Testing

  Testing is used to test coverage and validness of Cpick without requiring user input. Testing is useful for making sure Cpick can work on your device. Setting testing to `true` turns on the test mode, while setting testing to `false` keeps Cpick on normal mode.

All of the tests can be found in the [tests.go](https://github.com/ethanbaker/cpick/blob/master/tests.go) file.

#### Controls/Keys

If you are looking for a simplified version, check out the [docs](https://github.com/ethanbaker/cpick/blob/master/colors.go). Here is a more in-depth explanation.

You can also see a minified help version by pressing backtick (\`) while the Cpick application is running.

In Cpick, there are two main screens: the **hue screen** and the **saturation-value screen**.

The **hue screen** is the first screen that pops up. It contains a hue slider present at the top, two color tint previews on the left, and a table of different preset or custom colors at the bottom.

The **saturation-value screen** is the screen that pops up once you press enter to select a color on the **hue screen**. It contains one large color gradient in the center and two color tint previews on the right.

On the **hue screen**, you initially start off on the hue slider. This is where you can select a specific hue by navigating the table, either using h and l (from vim) or the left and right arrow keys. If you press enter, you will be moved to the **saturation-value** screen with the table representing the hue you chose. You can also press space to move to the listed colors at the bottom of the **hue screen**.

For the listed color table on the **hue screen**, you can navigate it using the standard vim keys (h, j, k, l) or the arrow keys. If you press enter on a selected color, you will be taken to the **saturation-value screen** and the color you selected will automatically be selected on the big color gradient table. You can also switch between color types by pressing c to go forward a page and C to go back.

In addition, you can press ? to open up a search menu. Here you can search for a specific color or keyword. You can press N to go to the next selection and n to go to the previous selection, just like in vim. 

On the **saturation-value screen**, you can move about the screen using the standard vim keys (h, j, k, l) or the arrow keys. Once on a desired color, press enter to select it (take note that if you ran Cpick with sandbox as true, pressing enter on the table will take you back to the **hue screen**).

In addition, you can press g to go to the start of the table (top-left most cell) or G to go to the end of the table (bottom-right most cell).

You can press 'q' or Escape to exit the application at any point, except in the search screen (Escape brings you back to the **hue screen** and 'q' is treated as a normal letter).

#### Custom Colors

In Cpick, you can add custom colors that can come up on the color pages. You can add JSON files that hold the colors in 3 ways.

1. **Local Environment**

Wherever you are running Cpick, you can provide a local `colors.json` file (file would have the path `./colors.json` from wherever Cpick is being run). This has the highest priority.

2. **~/.config Directory**

In the `~/.config` directory, you can create a `cpick` directory that can contain the `colors.json` file (file would have the path `~/.config/cpick/colors.json`). This has a lower priority than a local file but a higher priority than a `./cpick` directory.

3. **~/.cpick Directory**

In your home directory, you can create a `.cpick` directory that can contain the `colors.json` file (file would have the path `~/.cpick/colors.json`). This has the lowest priority.

The `colors.json` file has a very strict format. 

~~~json
{
  "colorList": [
    {
      "name": "COLOR NAME 1",
      "colors": []
    },

    {	
      "name": "COLOR NAME 2",
      "colors": []
    },

    {	
      "name": "COLOR NAME 3",
      "colors": []
    }

  ]
}
~~~

Each color type consists of an object with a `name` and `colors` key. The `name` key consists of what color type the colors provided are, such as CSS, Solarized, or XTERM. The `colors` key holds all of the colors that will be previewed when Cpick is run. All of the different color types are in the `colorList` array.

An individual color is an object that consists of two keys: `name` and `value`. `name` is the name of the color and `value` is the hexadecimal value of the color as a string. The "#" for the hexadecimal value is optional.

~~~json
"colors": [
  {"name": "Red",   "value": "#FF0000"},
  {"name": "Green", "value": "#00FF00"},
  {"name": "Blue",  "value": "#0000FF"}
]
~~~

Cpick comes with three color types as a default: CSS, Solarized, and XTERM. In order to fix complicated import problems, the JSON data is present in the [colors.go](https://github.com/ethanbaker/cpick/blob/master/colors.go) file as a string. The preset data always has the lowest priority for being used.


<p align="right">(<a href="#top">back to top</a>)</p>


<!-- ROADMAP -->
## Roadmap

- [x] Version 1.2.0
    - [x] Color Searching
- [ ] Version 1.3.0
    - [ ] Black Box Testing
    - [ ] Application Breakpoints

See the [open issues][issues-url] for a full list of proposed features (and known issues).

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTRIBUTING -->
## Contributing

For issues and suggestions, please include as much useful information as possible. Review the [documentation][documentation-url] and make sure the issue is actually present or the suggestion is not included. Please share issues/suggestions on the [issue tracker][issues-url].

For patches and feature additions, please submit them as [pull requests][pulls-url].  Please adhere to the [conventional commits][conventional-commits-url]. standard for commit messaging. In addition, please try to name your git branch according to your new patch. [These standards][conventional-branches-url] are a great guide you can follow.

You can follow these steps below to create a pull request:

1. Fork the Project
2. Create your Feature Branch (`git checkout -b branch_name`)
3. Commit your Changes (`git commit -m "commit_message"`)
4. Push to the Branch (`git push origin branch_name`)
5. Open a Pull Request

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- LICENSE -->
## License

This project uses the Apache 2.0 license.

You can find more information about this license in the `LICENSE` file.

<p align="right">(<a href="#top">back to top</a>)</p>


<!-- CONTACT -->
## Contact

Ethan Baker - contact@ethanbaker.dev - [LinkedIn][linkedin-url]

Project Link: [https://github.com/ethanbaker/cpick][project-url]

<p align="right">(<a href="#top">back to top</a>)</p>



<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/ethanbaker/cpick.svg
[forks-shield]: https://img.shields.io/github/forks/ethanbaker/cpick.svg
[stars-shield]: https://img.shields.io/github/stars/ethanbaker/cpick.svg
[issues-shield]: https://img.shields.io/github/issues/ethanbaker/cpick.svg
[license-shield]: https://img.shields.io/github/license/ethanbaker/cpick.svg
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?logo=linkedin&colorB=555

[contributors-url]: <https://github.com/ethanbaker/cpick/graphs/contributors>
[forks-url]: <https://github.com/ethanbaker/cpick/network/members>
[stars-url]: <https://github.com/ethanbaker/cpick/stargazers>
[issues-url]: <https://github.com/ethanbaker/cpick/issues>
[pulls-url]: <https://github.com/ethanbaker/cpick/pulls>
[license-url]: <https://github.com/ethanbaker/cpick/blob/master/LICENSE>
[linkedin-url]: <https://linkedin.com/in/ethandbaker>
[project-url]: <https://github.com/ethanbaker/cpick>

[product-screenshot]: ./docs/demonstration.gif
[documentation-url]: <https://pkg.go.dev/github.com/ethanbaker/cpick>

[conventional-commits-url]: <https://www.conventionalcommits.org/en/v1.0.0/#summary>
[conventional-branches-url]: <https://docs.microsoft.com/en-us/azure/devops/repos/git/git-branching-guidance?view=azure-devops>