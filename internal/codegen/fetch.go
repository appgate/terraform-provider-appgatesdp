package codegen

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apiversion"
	"gopkg.in/yaml.v3"
)

// Fetcher retrieves specifications from the sdp-api-specification repository
// with the git CLI. Extracted trees are cached per pinned commit, so repeated
// generation is offline and reproducible.
type Fetcher struct {
	Repository string
	CacheDir   string
}

func (f *Fetcher) git(ctx context.Context, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", append([]string{"--git-dir", f.repoDir()}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return out, fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(stderr.String()))
	}
	return out, nil
}

func (f *Fetcher) repoDir() string { return filepath.Join(f.CacheDir, "sdp-api-specification.git") }

func (f *Fetcher) ensureRepo(ctx context.Context) error {
	if _, err := os.Stat(f.repoDir()); err == nil {
		// The cache is disposable: recreate it if it is no longer a valid
		// repository (e.g. its files were removed by a cleanup tool).
		if _, err := f.git(ctx, "rev-parse", "--git-dir"); err != nil {
			if err := os.RemoveAll(f.repoDir()); err != nil {
				return err
			}
		}
	}
	if _, err := os.Stat(f.repoDir()); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(f.repoDir(), 0o755); err != nil {
			return err
		}
		if _, err := f.git(ctx, "init", "--bare", "--quiet"); err != nil {
			return err
		}
		_, err = f.git(ctx, "remote", "add", "origin", f.Repository)
		return err
	}
	_, err := f.git(ctx, "remote", "set-url", "origin", f.Repository)
	return err
}

// BranchHead returns the commit the branch currently points to.
func (f *Fetcher) BranchHead(ctx context.Context, branch string) (string, error) {
	if err := f.ensureRepo(ctx); err != nil {
		return "", err
	}
	out, err := f.git(ctx, "ls-remote", "--heads", "origin", "refs/heads/"+branch)
	if err != nil {
		return "", err
	}
	sha, _, _ := strings.Cut(strings.TrimSpace(string(out)), "\t")
	if sha == "" {
		return "", fmt.Errorf("branch %s not found in %s", branch, f.Repository)
	}
	return sha, nil
}

// Fetch extracts the specification of v at its pinned commit and returns the
// directory. It fails when the branch is missing, the commit is not on the
// branch, or the tree does not contain the expected API version.
func (f *Fetcher) Fetch(ctx context.Context, v apiversion.Version) (string, error) {
	dir := filepath.Join(f.CacheDir, "specs", v.APIVersion, v.Commit)
	if _, err := os.Stat(filepath.Join(dir, ".complete")); err == nil {
		return dir, validateSpecTree(dir, v)
	}
	if _, err := f.BranchHead(ctx, v.Branch); err != nil {
		return "", fmt.Errorf("%s: %w", v.APIVersion, err)
	}
	ref := "refs/remotes/origin/" + v.Branch
	if _, err := f.git(ctx, "fetch", "--quiet", "--no-tags", "origin", "+refs/heads/"+v.Branch+":"+ref); err != nil {
		return "", fmt.Errorf("%s: fetch branch %s: %w", v.APIVersion, v.Branch, err)
	}
	if _, err := f.git(ctx, "cat-file", "-e", v.Commit+"^{commit}"); err != nil {
		return "", fmt.Errorf("%s: pinned commit %s is not part of branch %s", v.APIVersion, v.Commit, v.Branch)
	}
	if _, err := f.git(ctx, "merge-base", "--is-ancestor", v.Commit, ref); err != nil {
		return "", fmt.Errorf("%s: pinned commit %s is not part of branch %s", v.APIVersion, v.Commit, v.Branch)
	}
	archive, err := f.git(ctx, "archive", "--format=tar", v.Commit)
	if err != nil {
		return "", err
	}
	tmp := dir + ".tmp"
	_ = os.RemoveAll(tmp)
	if err := untar(bytes.NewReader(archive), tmp); err != nil {
		return "", err
	}
	if err := validateSpecTree(tmp, v); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(tmp, ".complete"), []byte(v.Commit+"\n"), 0o644); err != nil {
		return "", err
	}
	_ = os.RemoveAll(dir)
	return dir, os.Rename(tmp, dir)
}

var specVersionRE = regexp.MustCompile(`API version (\d+)`)

// validateSpecTree checks that dir holds the root specification of v.
func validateSpecTree(dir string, v apiversion.Version) error {
	b, err := os.ReadFile(filepath.Join(dir, "api_specs.yml"))
	if err != nil {
		return fmt.Errorf("%s: specification root api_specs.yml missing at commit %s: %w", v.APIVersion, v.Commit, err)
	}
	var doc struct {
		OpenAPI string `yaml:"openapi"`
		Info    struct {
			Version string `yaml:"version"`
		} `yaml:"info"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return fmt.Errorf("%s: invalid api_specs.yml: %w", v.APIVersion, err)
	}
	if !strings.HasPrefix(doc.OpenAPI, "3.") {
		return fmt.Errorf("%s: api_specs.yml is not an OpenAPI 3 document (openapi: %q)", v.APIVersion, doc.OpenAPI)
	}
	m := specVersionRE.FindStringSubmatch(doc.Info.Version)
	if m == nil || m[1] != fmt.Sprint(v.Number()) {
		return fmt.Errorf("%s: specification at %s declares %q, expected API version %d", v.APIVersion, v.Commit, doc.Info.Version, v.Number())
	}
	return nil
}

func untar(r io.Reader, dst string) error {
	tr := tar.NewReader(r)
	for {
		h, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		name := filepath.Clean(h.Name)
		if filepath.IsAbs(name) || strings.HasPrefix(name, "..") {
			return fmt.Errorf("unsafe path %q in archive", h.Name)
		}
		target := filepath.Join(dst, name)
		switch h.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			b, err := io.ReadAll(tr)
			if err != nil {
				return err
			}
			if err := os.WriteFile(target, b, 0o644); err != nil {
				return err
			}
		}
	}
}
