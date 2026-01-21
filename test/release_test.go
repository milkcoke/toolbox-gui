package test

import (
	"errors"
	"fyne.io/fyne/v2"
	"golang.org/x/sys/execabs"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
	"text/template/parse"
)

type FuncMap map[string]any

type missingKeyAction int

const (
	mapInvalid   missingKeyAction = iota // Return an invalid reflect.Value.
	mapZeroValue                         // Return the zero value for the map element.
	mapError                             // Error out
)

type option struct {
	missingKey missingKeyAction
}
type common struct {
	tmpl   map[string]*Template // Map from name to defined templates.
	muTmpl sync.RWMutex         // protects tmpl
	option option
	// We use two maps, one for parsing and one for execution.
	// This separation makes the API cleaner since it doesn't
	// expose reflection to the client.
	muFuncs    sync.RWMutex // protects parseFuncs and execFuncs
	parseFuncs FuncMap
	execFuncs  map[string]reflect.Value
}

type Template struct {
	name string
	*parse.Tree
	*common
	leftDelim  string
	rightDelim string
}

type keyValueFlag struct {
	m map[string]string
}

type appData struct {
	icon, Name        string
	AppID, AppVersion string
	AppBuild          int
	ResGoString       string
	Release           bool
	CustomMetadata    map[string]string
	VersionAtLeast2_3 bool
}

var resourceEntitlementsDarwinPlist = &fyne.StaticResource{
	StaticName: "entitlements-darwin.plist",
	StaticContent: []byte(
		"<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n<!DOCTYPE plist PUBLIC \"-//Apple//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">\n<plist version=\"1.0\">\n<dict>\n    <key>com.apple.security.app-sandbox</key>\n    <true/>\n</dict>\n</plist>"),
}
var EntitlementsDarwin = template.Must(template.New("Entitlements").Parse(string(resourceEntitlementsDarwinPlist.StaticContent)))

type Packager struct {
	*appData
	srcDir, dir, exe, os           string
	install, release, distribution bool
	certificate, profile           string // optional flags for releasing
	tags, category                 string
	tempDir                        string

	customMetadata keyValueFlag
}

type Releaser struct {
	Packager

	keyStore     string
	keyStorePass string
	keyName      string
	keyPass      string
	developer    string
	password     string
}

func (r *Releaser) writeEntitlements(tmpl *template.Template, entitlementData interface{}) (cleanup func(), err error) {
	entitlementPath := filepath.Join(r.dir, "entitlements.plist")
	entitlements, err := os.Create(entitlementPath)
	log.Println("path entitle : ", entitlementPath)
	log.Println("entitlements: ", entitlements)
	if err != nil {
		return nil, err
	}
	defer func() {
		if r := entitlements.Close(); r != nil && err == nil {
			err = r
		}
	}()

	if err := tmpl.Execute(entitlements, entitlementData); err != nil {
		return nil, err
	}
	return nil, nil
	//return func() {
	//	_ = os.Remove(entitlementPath)
	//}, nil
}

func (r *Releaser) packageMacOSRelease() error {
	// try to derive two certificates from one name (they will be consistent)
	appCert := strings.Replace(r.certificate, "Installer", "Application", 1)
	log.Println("appCert "+": ", appCert)
	installCert := strings.Replace(r.certificate, "Application", "Installer", 1)
	log.Println("installCert "+": ", installCert)
	unsignedPath := r.Name + "-unsigned.pkg"

	//defer os.RemoveAll(r.Name + ".app") // this was the output of package and it can get in the way of future builds

	_, err := r.writeEntitlements(EntitlementsDarwin, nil)
	if err != nil {
		return errors.New("failed to write entitlements plist template")
	}
	//defer cleanup()

	//# sign the app:
	//codesign --force --options runtime --deep --sign "${CERT}" -i "${APP_ID}" "${APP_NAME}.app" || exit

	cmd := execabs.Command("codesign", "-vfs", appCert, "--entitlement", "entitlements.plist", r.Name+".app")
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Codesign failed", err)
		return errors.New("unable to codesign application bundle")
	}

	cmd = execabs.Command("productbuild", "--component", r.Name+".app", "/Applications/",
		"--product", r.Name+".app/Contents/Info.plist", unsignedPath)
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Product build failed", err)
		return errors.New("unable to build macOS app package")
	}
	//defer os.Remove(unsignedPath)

	cmd = execabs.Command("productsign", "--sign", installCert, unsignedPath, r.Name+".pkg")
	return cmd.Run()
}

func Test_MakePkg(t *testing.T) {
	curDir, err := os.Getwd()

	rootDir := filepath.Join(curDir, "..", "..")

	appData := &appData{
		icon:       filepath.Join(rootDir, "./assets/app_icon.png"),
		AppID:      "B3VZMK576W.toolbox.milkcoke",
		AppVersion: "1.0.0",
		Name:       "Toolbox Dev",
		AppBuild:   1,
		Release:    true,
	}
	//fyne release -icon assets/app_icon.png -name "Toolbox Dev" \
	//--os darwin -appID "B3VZMK576W.toolbox.milkcoke" --appVersion 1.0 --appBuild 1 \
	//--cert "3rd Party Mac Developer Application: SeungHun Moon (B3VZMK576W)" \
	//--profile "toolbox_milkcoke_profile.provisionprofile" \
	//-category "developer-tools"

	packageDir := filepath.Join(rootDir, "test_package")

	packager := Packager{
		appData:     appData,
		srcDir:      rootDir,
		dir:         packageDir,
		os:          "darwin",
		tempDir:     "test_package",
		release:     true,
		certificate: "3rd Party Mac Developer Application: SeungHun Moon (B3VZMK576W)",
		category:    "developer-tools",
	}
	r := Releaser{
		Packager: packager,
	}

	err = r.packageMacOSRelease()
	if err != nil {
		log.Fatalln(err)
		return
	}
}
