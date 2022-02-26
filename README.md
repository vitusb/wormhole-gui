# Wormhole Gui [(deutsche Beschreibung)](https://github.com/vitusb/wormhole-gui/blob/main/wormhole-gui.txt) 🇩🇪

"Wormhole Gui" ist eine plattformübergreifende Anwendung, mit der Sie Dateien, Ordner und Text problemlos zwischen Geräten austauschen können. Es verwendet die Go-Implementierung von magic-wormhole, genannt [wormhole-william](https://github.com/vitusb/wormhole-william), und wird statisch in eine einzige Binärdatei kompiliert. Wormhole-gui ist auch mit dem Senden an und Empfangen von anderen Wormhole-Clients kompatibel, wie z. B. den CLI-Anwendungen von [wormhole-william](https://github.com/vitusb/wormhole-william) und [magic-wormhole](https://github.com/magic-wormhole/magic-wormhole). 

<p align="center">
  <img src="internal/assets/screenshot.png" />
</p>

Erstellt unter Verwendung folgender Go-Lang Module:
- [fyne](https://github.com/fyne-io/fyne) (version 2.1.3 / german patches)
- [wormhole-william](https://github.com/vitusb/wormhole-william) (latest github-master)
- [compress](https://github.com/klauspost/compress) (version 1.14.4)

# wormhole-gui (German Version) 🇩🇪

This is the first attempt to create a german version of Wormhole-Gui inclusive german PGP-wordlist (by german version of wormhole-william). Due to the fact that Fyne has no native language-support, i'm forced to translate Fyne [directly inline in code](https://github.com/vitusb/fyne/tree/master/german). I hope that this will not be a challenge for future code-releases. I will do my very best ...

        - vitusb- 😄

## - English description - 🇬🇧

The initial version was built in less than one day to show how quick and easy it is to use [fyne](https://github.com/fyne-io/fyne) for developing applications.

## Sponsoring

Wormhole-gui an open source project that is provided free of charge and that will continue to be the case forever. If you use wormhole-gui and appreciate the work being put into it, please consider supporting the development through [GitHub Sponsors](https://github.com/sponsors/Jacalz). This is in no way a requirement, but would be greatly appreciated and would allow for even more improvements to come further down the road.

## Requirements

Wormhole-gui compiles into a statically linked binary with no runtime dependencies.
Compiling requires a [Go](https://golang.org) compiler (1.14 or later) and the [prerequisites for Fyne](https://developer.fyne.io/started/).

## Downloads

Please visit the [release page](https://github.com/vitusb/wormhole-gui/releases) for downloading the latest releases.
Pre-built binaries are available for FreeBSD, Linux, MacOS (`x86-64` and `arm64`) and Windows (`x86-64`).

The following distrubutions also have binary packages avaliable through the respective package managers:

[![Packaging status](https://repology.org/badge/vertical-allrepos/wormhole-gui.svg)](https://repology.org/project/wormhole-gui/versions)

## Building

Systems with the compile-time requirements satisfied can build the project using `go build` in the project root:
```bash
go build
```

The project is available in the [Fyne Apps Listing](https://apps.fyne.io/apps/wormhole-gui.html) and can be installed either using the `fyne get` command or using the [Fyne Apps Installer](https://apps.fyne.io/apps/io.fyne.apps.html).
It can also be built and installed using GNU Make (installing is currently only supported on Linux and BSD):
```bash
make
sudo make install
```

## Contributing

Contributions are strongly appreciated. Everything from creating bug reports to contributing code will help the project a lot, so please feel free to help in any way, shape or form that you feel comfortable doing.

## License
- Wormhole-gui is licensed under `GNU GENERAL PUBLIC LICENSE Version 3`.
