package dockerfile

import (
	"io/ioutil"
	"testing"
	"os"
	
	"gotest.tools/assert"
)

func TestGetPorts(t *testing.T) {
	dir, err := ioutil.TempDir("", "testDeploy")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}

	wdBackup, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		t.Fatalf("Error changing working directory: %v", err)
	}

	// 8. Delete temp folder
	defer func() {
		err = os.Chdir(wdBackup)
		if err != nil {
			t.Fatalf("Error changing dir back: %v", err)
		}
		err = os.RemoveAll(dir)
		if err != nil {
			t.Fatalf("Error removing dir: %v", err)
		}
	}()

	file, err := os.Create("Dockerfile")
	if err != nil {
		t.Fatalf("Error creating Dockerfile: %v", err)
	}
	_, err = file.Write([]byte(`FROM mysql
EXPOSE 8080
EXPOSE `))
	if err != nil {
		t.Fatalf("Error creating Dockerfile: %v", err)
	}
	err = file.Close()
	if err != nil {
		t.Fatalf("Error creating Dockerfile: %v", err)
	}

	ports, err := GetPorts("Dockerfile")
	if err != nil {
		t.Fatalf("Error receiving ports: %v", err)
	}
	assert.Equal(t, 1, len(ports), "Wrong number of ports returned")
	assert.Equal(t, 8080, ports[0], "Wrong port returned")


}
