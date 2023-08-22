

<p align="center">
  <img src="https://github.com/joseviniciusnunes/sh-icon-tray/assets/22475804/0404a2cc-3091-4497-a365-e750a9bb1c01" height="300px" />
  <h1 align="center">sh-icon-tray</h1>
</p>

## Create quick shortcuts in the system tray using bash or cmd.


#### 1º: Download the binary for your platform (MacOS, Linux and Windows) in the releases area.
#### 2º: Place the binary to start next to your operating system.
#### 3º: Run the binary for the first time, it will create a file called sh-icon-tray.yml in your user's HOME, you can edit it by clicking on More -> Edit Config, customize with your commands.

## sh-tray-icon.yml example

```yaml
root:
  - label: My Projects
    children:
      - label: my-app-react
        run: code ~/projects/my-app-react

      - label: my-app-golang
        run: code ~/projects/my-app-golang

  - divider: true

  - label: SYSTEM   docker stats
    run: docker stats

  - label: SYSTEM   htop
    run: htop

```
