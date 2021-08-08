# dots-cli

`dots` is CLI tool to build, version and publish config file bundles.

## TODO (shared):
- [ ] Writing comprehensive documentation
- [ ] Designing landing page
- [ ] Providing screenshots/videos for demonstration

## TODO (client):
### Functionality:
- [ ] Commands:
    - [x] cmd: `init` (initializes empty package)
    - [x] cmd: `add` (adds app to package)
    - [x] cmd: `remove` (removes app from package)
    - [x] cmd: `pack` (saves package's current state)
    - [x] cmd: `login` (signs user in and saves token)
    - [x] cmd: `push` (pushes package to added registry)
    - [x] cmd: `remote` (adds/removes registries)
      - [x] `add` (adds remote registry)
      - [x] `remove` (removes remote registry)
    - [x] cmd: `list`(lists apps)
      - [x] `added` (lists added apps in package)
      - [x] `all` (use --installed flag to list only apps found on system) (lists all possible apps from config)
    - [x] cmd: `install` (installs packages to appropriate directories)
        - [x] Backup before installation
        - [ ] Revert if installation failes
    - [x] cmd: `get` (downloads package)
    - [x] cmd: Find better place for logs.txt + better logging
    - [ ] cmd: `revert` (loads previous version from package archive)

## TODO ([dots-server](github.com/alvanrahimli/dots-server))
### Functionality:
- [x] Modelling
- [x] Login / Register API endpoints
- [x] Push package API endpoint
- [x] Get PackageArchive endpoint
- [ ] Update & Delete packages endpoint
- [ ] Enhanced models (info, settings etc.)
- [ ] Enhanced endpoints for webapp
- [ ] Security considerations
