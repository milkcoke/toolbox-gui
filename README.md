

#### Prerequisite
- [GNU make](https://gnuwin32.sourceforge.net/packages/make.htm)
- [Go](https://go.dev/dl/)
- [Fyne](https://developer.fyne.io/started/) 
- msys2 \
C compiler for windows

Enable Go Modules integration in IntelliJ IDEA.

## What's the Fyne?
Fyne uses single binding to compile our app with existing graphics driver independent on Operation System.

### Wrapper
This is used to structure lay out

## What's learn from this
- [X] How to get widget
- [ ] How to generate Image button (Toggle)
- [ ] How to generate check box 
- [X] How to package our application
- Button Event Listener 
  - [ ] Batch Download
  - [ ] Download one by one
  - [ ] Setting download target path
  - [ ] Open directory

### Packaging
#### 1. Install fyne cmd
```bash
$ go install fyne.io/fyne/v2/cmd/fyne@latest
```

#### 2. Prepare icon of app

#### 3. Input cli
Refer to [this docs](https://developer.fyne.io/started/packaging)
```bash
# change directory including main.go
$ cd [PATH]
$ fyne package --appVersion [version] -name [app-name] -icon [app-icon-path] -release
```

> This is just only usable on your machine. \
> This binary files generated by cli doesn't work any other OS.

### Register application
You need give some money about $100 / year to Apple or Microsoft for supporting cross-platform \
for signing your application assuring your app is secure.


## Test
I defined `Makefile` for integration testing.
```bash
# You can execute all test codes in this project.
$ make test
```