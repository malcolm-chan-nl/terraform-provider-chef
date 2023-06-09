package provider

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// To run these acceptance tests, you will need access to a Chef server.
// An easy way to get one is to sign up for a hosted Chef server account
// at https://manage.chef.io/signup , after which your base URL will
// be something like https://api.opscode.com/organizations/example/ .
// You will also need to create a "client" and write its private key to
// a file somewhere.
//
// You can then set the following environment variables to make these
// tests work:
// CHEF_SERVER_URL to the base URL as described above.
// CHEF_CLIENT_NAME to the name of the client object you created.
// CHEF_KEY_MATERIAL the key file contents.
//
// You will probably need to edit the global permissions on your Chef
// Server account to allow this client (or all clients, if you're lazy)
// to have both List and Create access on all types of object:
//     https://manage.chef.io/organizations/yourorg/global_permissions
//
// With all of that done, you can run like this:
//    make testacc TEST=./builtin/providers/chef

var testAccProvider *schema.Provider

func init() {
	testAccProvider = New("dev")()
	if testSuffix == "" {
		bytes := make([]byte, 4)
		if _, err := rand.Read(bytes); err != nil {
			log.Fatal(err)
		}
		testSuffix = hex.EncodeToString(bytes)
	}
}

var testSuffix string

func testSuffixRender(cfg string) string {
	t := template.Must(template.New("").Parse(cfg))
	var b bytes.Buffer
	t.Execute(&b, testSuffix)

	return b.String()
}

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"chef": func() (*schema.Provider, error) {
		return testAccProvider, nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CHEF_SERVER_URL"); v == "" {
		t.Fatal("CHEF_SERVER_URL must be set for acceptance tests")
	}
	if v := os.Getenv("CHEF_CLIENT_NAME"); v == "" {
		t.Fatal("CHEF_CLIENT_NAME must be set for acceptance tests")
	}
	if v := os.Getenv("CHEF_KEY_MATERIAL"); v == "" {
		t.Fatal("CHEF_KEY_MATERIAL must be set for acceptance tests")
	}
}
