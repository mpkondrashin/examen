package update

import "testing"

const repo = "sandboxer"

func TestCheckUpdate(t *testing.T) {
	u, err := CheckLocationGithub(repo)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("url = %v", u)
	version := ParseVersion(u)
	t.Logf("version = %v", version)
	err = DownloadRelease(version, "Setup.zip", ".", func(p float32) error {
		t.Logf("Downloaded %d%%", int(p*100))
		return nil
	})
	if err != nil {
		t.Error(err)
	}

}
