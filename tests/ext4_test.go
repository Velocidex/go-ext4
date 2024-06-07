package tests

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Velocidex/ordereddict"
	"github.com/alecthomas/assert"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
)

type EXT4TestSuite struct {
	suite.Suite
	binary, extension string
	tmpdir            string
}

func (self *EXT4TestSuite) SetupTest() {
	if runtime.GOOS == "windows" {
		self.extension = ".exe"
	}

	// Search for a valid binary to run.
	binaries, err := filepath.Glob(
		"../goext4" + self.extension)
	assert.NoError(self.T(), err)

	self.binary, _ = filepath.Abs(binaries[0])
	fmt.Printf("Found binary %v\n", self.binary)

	self.tmpdir, err = ioutil.TempDir("", "tmp")
	assert.NoError(self.T(), err)
}

func (self *EXT4TestSuite) TearDownTest() {
	os.RemoveAll(self.tmpdir)
}

func (self *EXT4TestSuite) TestExt4() {
	image := "./images/hierarchy_64.ext4"

	cmd := exec.Command(self.binary, "ls", "-r", image, "/")
	out, err := cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden := ordereddict.NewDict().Set("listdir", strings.Split(string(out), "\n"))

	cmd = exec.Command(self.binary, "cat", image,
		"/directory1/subdirectory1/fortune3")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	md5_sum := md5.New()
	md5_sum.Write(out)

	golden.Set("fortune3_md5", hex.EncodeToString(md5_sum.Sum(nil)))

	cmd = exec.Command(self.binary, "stat", image, "thejungle.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden.Set("stat", strings.Split(string(out), "\n"))

	g := goldie.New(self.T(), goldie.WithFixtureDir("fixtures"))
	g.AssertJson(self.T(), "hierarchy_64.ext4", golden)
}

func TestExt4(t *testing.T) {
	suite.Run(t, &EXT4TestSuite{})
}

func formatForTest(v interface{}) []string {
	serialized, _ := json.MarshalIndent(v, " ", " ")
	return strings.Split(string(serialized), "\n")
}
